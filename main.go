package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/alertmanager"
	"github.com/thomasobenaus/sokar/api"
	"github.com/thomasobenaus/sokar/capacityPlanner"
	"github.com/thomasobenaus/sokar/config"
	"github.com/thomasobenaus/sokar/logging"
	"github.com/thomasobenaus/sokar/nomad"
	"github.com/thomasobenaus/sokar/scaleAlertAggregator"
	"github.com/thomasobenaus/sokar/scaler"
	"github.com/thomasobenaus/sokar/sokar"
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// cliAndConfig provides the configuration by reading parameters from the cli and from config-file.
func cliAndConfig(args []string) (*config.Config, error) {

	// parse commandline args and consume environment variables
	parsedArgs, err := parseArgs(args)
	if err != nil {
		return nil, err
	}

	if !parsedArgs.validateArgs() {
		parsedArgs.printDefaults()
		return nil, fmt.Errorf("Invalid cli parameters")
	}

	log.Println("Read configuration...")
	cfg, err := config.NewConfigFromYAMLFile(parsedArgs.CfgFile)
	if err != nil {
		log.Printf("Error reading configuration: %s. Using the default config instead.", err.Error())
	}

	// Prefer CLI parameter for the nomadServerAddress
	if len(parsedArgs.NomadServerAddr) > 0 {
		cfg.Nomad.ServerAddr = parsedArgs.NomadServerAddr
	}

	if len(cfg.Nomad.ServerAddr) == 0 {
		parsedArgs.printDefaults()
		return nil, fmt.Errorf("Nomad Server address not specified")
	}

	log.Println("Read configuration...done")
	return &cfg, nil
}

// setupLogging configures logging according to the given parameters
func setupLogging(cfg *config.Config) (logging.LoggerFactory, error) {
	if cfg == nil {
		return nil, fmt.Errorf("Error creating LoggerFactory: Config is nil")
	}
	lCfg := logging.Config{
		UseStructuredLogging:       cfg.Logging.Structured,
		UseUnixTimestampForLogging: cfg.Logging.UxTimestamp,
	}
	loggingFactory := lCfg.New()
	return loggingFactory, nil
}

func setupScaleAlertAggregator(scaleAlertEmitters []scaleAlertAggregator.ScaleAlertEmitter, cfg *config.Config, logF logging.LoggerFactory) *scaleAlertAggregator.ScaleAlertAggregator {
	weightMap := make(scaleAlertAggregator.ScaleAlertWeightMap, 0)
	for _, alertDef := range cfg.ScaleAlertAggregator.ScaleAlerts {
		weightMap[alertDef.Name] = alertDef.Weight
	}

	scaEvtAggCfg := scaleAlertAggregator.Config{
		Logger:                 logF.NewNamedLogger("sokar.scaAlertAggr"),
		NoAlertScaleDamping:    cfg.ScaleAlertAggregator.NoAlertScaleDamping,
		UpScalingThreshold:     cfg.ScaleAlertAggregator.UpScaleThreshold,
		DownScalingThreshold:   cfg.ScaleAlertAggregator.DownScaleThreshold,
		EvaluationCycle:        cfg.ScaleAlertAggregator.EvaluationCycle,
		EvaluationPeriodFactor: cfg.ScaleAlertAggregator.EvaluationPeriodFactor,
		CleanupCycle:           cfg.ScaleAlertAggregator.CleanupCycle,
		WeightMap:              weightMap,
	}

	scaAlertAggr := scaEvtAggCfg.New(scaleAlertEmitters, scaleAlertAggregator.NewMetrics())
	return scaAlertAggr
}

func setupScaleAlertEmitters(api *api.API, logF logging.LoggerFactory) []scaleAlertAggregator.ScaleAlertEmitter {
	// Alertmanger Connector
	amCfg := alertmanager.Config{
		Logger: logF.NewNamedLogger("sokar.alertmanager"),
	}
	amConnector := amCfg.New()
	api.Router.POST("/alerts", amConnector.HandleScaleAlerts)

	var scaleAlertEmitters []scaleAlertAggregator.ScaleAlertEmitter
	scaleAlertEmitters = append(scaleAlertEmitters, amConnector)

	return scaleAlertEmitters
}

func main() {

	// read config
	cfg, err := cliAndConfig(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// set up logging
	loggingFactory, err := setupLogging(cfg)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	logger := loggingFactory.NewNamedLogger("sokar")
	logger.Info().Msg("Connecting components and setting up sokar ...")

	// 1. API
	api := api.New(cfg.Port, loggingFactory.NewNamedLogger("sokar.api"))

	// 2. AlertEmitters (i.e. Alertmanager Connector)
	scaleAlertEmitters := setupScaleAlertEmitters(api, loggingFactory)

	// 3. ScaleAlertAggregator
	scaAlertAggr := setupScaleAlertAggregator(scaleAlertEmitters, cfg, loggingFactory)

	// 4. Scaler
	scaler, err := setupScaler(cfg.Job.Name, cfg.Job.MinCount, cfg.Job.MaxCount, cfg.Nomad.ServerAddr, loggingFactory)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed setting up the scaler")
	}

	// 5. CapacityPlanner
	capaCfg := capacityPlanner.Config{
		Logger: loggingFactory.NewNamedLogger("sokar.capaPlanner"),
	}
	capaPlanner := capaCfg.New()

	// 6. Sokar
	sokarInst, err := setupSokar(scaAlertAggr, capaPlanner, scaler, api, logger)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed creating sokar.")
	}

	// Register metrics handler
	api.Router.Handler("GET", "/metrics", promhttp.Handler())

	logger.Info().Msg("Connecting components and setting up sokar ... done")

	// Define runnables and their execution order
	var orderedRunnables []Runnable
	orderedRunnables = append(orderedRunnables, sokarInst)
	orderedRunnables = append(orderedRunnables, scaler)
	orderedRunnables = append(orderedRunnables, scaAlertAggr)
	orderedRunnables = append(orderedRunnables, api)

	// Run all components
	Run(orderedRunnables, logger)

	// Install signal handler for shutdown
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, syscall.SIGINT, syscall.SIGTERM)
	go shutdownHandler(shutDownChan, orderedRunnables, logger)

	// Wait till completion
	Join(orderedRunnables, logger)

	logger.Info().Msg("Shutdown successfully completed")
	os.Exit(0)
}

func setupSokar(scaleEventEmitter sokarIF.ScaleEventEmitter, capacityPlanner sokarIF.CapacityPlanner, scaler sokarIF.Scaler, api *api.API, logger zerolog.Logger) (*sokar.Sokar, error) {
	cfg := sokar.Config{
		Logger: logger,
	}
	sokarInst, err := cfg.New(scaleEventEmitter, capacityPlanner, scaler)
	if err != nil {
		return nil, err
	}

	logger.Info().Msg("Registering http handlers ...")
	api.Router.POST(sokar.PathScaleBy, sokarInst.ScaleBy)
	api.Router.GET(sokar.PathHealth, sokarInst.Health)
	logger.Info().Msg("Registering http handlers ... done")

	return sokarInst, err
}

// setupScaler creates and configures the Scaler. Internally nomad is used as scaling target.
func setupScaler(jobName string, min uint, max uint, nomadSrvAddr string, logF logging.LoggerFactory) (*scaler.Scaler, error) {

	if logF == nil {
		return nil, fmt.Errorf("Logging factory is nil")
	}

	// Set up the nomad connector
	nomadConfig := nomad.NewDefaultConfig(nomadSrvAddr)
	nomadConfig.Logger = logF.NewNamedLogger("sokar.nomad")
	nomad, err := nomadConfig.New()
	if err != nil {
		return nil, fmt.Errorf("Failed setting up nomad connector: %s", err)
	}

	scaCfg := scaler.Config{
		JobName:  jobName,
		MinCount: min,
		MaxCount: max,
		Logger:   logF.NewNamedLogger("sokar.scaler"),
	}

	scaler, err := scaCfg.New(nomad)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up scaler: %s", err)
	}

	return scaler, nil
}
