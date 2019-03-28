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
	"github.com/thomasobenaus/sokar/helper"
	"github.com/thomasobenaus/sokar/logging"
	"github.com/thomasobenaus/sokar/nomad"
	"github.com/thomasobenaus/sokar/scaleAlertAggregator"
	"github.com/thomasobenaus/sokar/scaler"
	"github.com/thomasobenaus/sokar/sokar"
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// read config
	cfg := helper.Must(cliAndConfig(os.Args)).(*config.Config)

	// set up logging
	loggingFactory := helper.Must(setupLogging(cfg)).(logging.LoggerFactory)

	logger := loggingFactory.NewNamedLogger("sokar")
	logger.Info().Msg("Connecting components and setting up sokar")

	logger.Info().Msg("1. Setup: API")
	api := api.New(cfg.Port, loggingFactory.NewNamedLogger("sokar.api"))

	logger.Info().Msg("2. Setup: ScaleAlertEmitters")
	scaleAlertEmitters := helper.Must(setupScaleAlertEmitters(api, loggingFactory)).([]scaleAlertAggregator.ScaleAlertEmitter)

	logger.Info().Msg("3. Setup: ScaleAlertAggregator")
	scaAlertAggr := setupScaleAlertAggregator(scaleAlertEmitters, cfg, loggingFactory)

	logger.Info().Msg("4. Setup: Scaler")
	scaler := helper.Must(setupScaler(cfg.Job.Name, cfg.Job.MinCount, cfg.Job.MaxCount, cfg.Nomad.ServerAddr, loggingFactory)).(*scaler.Scaler)

	logger.Info().Msg("5. Setup: CapacityPlanner")
	capaCfg := capacityPlanner.NewDefaultConfig()
	capaCfg.Logger = loggingFactory.NewNamedLogger("sokar.capaPlanner")
	capaPlanner := capaCfg.New()

	logger.Info().Msg("6. Setup: Sokar")
	sokarInst := helper.Must(setupSokar(scaAlertAggr, capaPlanner, scaler, api, logger)).(*sokar.Sokar)

	// Register metrics handler
	api.Router.Handler("GET", sokar.PathMetrics, promhttp.Handler())

	// Define runnables and their execution order
	var orderedRunnables []Runnable
	orderedRunnables = append(orderedRunnables, sokarInst)
	orderedRunnables = append(orderedRunnables, scaler)
	orderedRunnables = append(orderedRunnables, scaAlertAggr)
	orderedRunnables = append(orderedRunnables, api)

	// Install signal handler for shutdown
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, syscall.SIGINT, syscall.SIGTERM)
	go shutdownHandler(shutDownChan, orderedRunnables, logger)

	// Run all components
	Run(orderedRunnables, logger)

	// Wait till completion
	Join(orderedRunnables, logger)

	logger.Info().Msg("Shutdown successfully completed")
	os.Exit(0)
}

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

func setupScaleAlertEmitters(api *api.API, logF logging.LoggerFactory) ([]scaleAlertAggregator.ScaleAlertEmitter, error) {
	if api == nil {
		return nil, fmt.Errorf("API is nil")
	}

	if logF == nil {
		return nil, fmt.Errorf("LoggingFactory is nil")
	}

	// Alertmanger Connector
	amCfg := alertmanager.Config{
		Logger: logF.NewNamedLogger("sokar.alertmanager"),
	}
	amConnector := amCfg.New()
	api.Router.POST(sokar.PathAlertmanager, amConnector.HandleScaleAlerts)

	var scaleAlertEmitters []scaleAlertAggregator.ScaleAlertEmitter
	scaleAlertEmitters = append(scaleAlertEmitters, amConnector)

	return scaleAlertEmitters, nil
}

func setupSokar(scaleEventEmitter sokarIF.ScaleEventEmitter, capacityPlanner sokarIF.CapacityPlanner, scaler sokarIF.Scaler, api *api.API, logger zerolog.Logger) (*sokar.Sokar, error) {
	cfg := sokar.Config{
		Logger: logger,
	}
	sokarInst, err := cfg.New(scaleEventEmitter, capacityPlanner, scaler, sokar.NewMetrics())
	if err != nil {
		return nil, err
	}

	api.Router.GET(sokar.PathHealth, sokarInst.Health)

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

	scaler, err := scaCfg.New(nomad, scaler.NewMetrics())
	if err != nil {
		return nil, fmt.Errorf("Failed setting up scaler: %s", err)
	}

	return scaler, nil
}
