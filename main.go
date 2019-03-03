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
)

func main() {

	// parse commandline args and consume environment variables
	parsedArgs := parseArgs()
	if !parsedArgs.validateArgs() {
		os.Exit(1)
	}

	localPort := 11000

	log.Println("Read configuration...")
	cfg, err := config.NewConfigFromYAMLFile(parsedArgs.CfgFile)
	if err != nil {
		log.Fatalf("Error reading configuration: %s.", err.Error())
	}
	log.Println("Read configuration...done")

	// set up logging
	lCfg := logging.Config{
		UseStructuredLogging:       cfg.Logging.Structured,
		UseUnixTimestampForLogging: cfg.Logging.UxTimestamp,
	}
	loggingFactory := lCfg.New()
	logger := loggingFactory.NewNamedLogger("sokar")

	nomadServerAddress := cfg.Nomad.ServerAddr
	// Prefer CLI parameter
	if len(parsedArgs.NomadServerAddr) > 0 {
		nomadServerAddress = parsedArgs.NomadServerAddr
	}
	if len(nomadServerAddress) == 0 {
		logger.Fatal().Msg("Nomad Server address not specified.")
	}

	logger.Info().Msg("Set up the scaler ...")
	scaler, err := setupScaler(cfg.Job.Name, cfg.Job.MinCount, cfg.Job.MaxCount, nomadServerAddress, loggingFactory)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed setting up the scaler")
	}
	logger.Info().Msg("Set up the scaler ... done")

	logger.Info().Msg("Connecting components and setting up sokar ...")
	api := api.New(localPort, loggingFactory.NewNamedLogger("sokar.api"))

	var scaleAlertEmitters []scaleAlertAggregator.ScaleAlertEmitter
	amCfg := alertmanager.Config{
		Logger: loggingFactory.NewNamedLogger("sokar.alertmanager"),
	}
	amConnector := amCfg.New()
	api.Router.POST("/alerts", amConnector.HandleScaleAlerts)
	scaleAlertEmitters = append(scaleAlertEmitters, amConnector)

	weightMap := make(scaleAlertAggregator.ScaleAlertWeightMap, 0)
	for _, alertDef := range cfg.ScaleAlertAggregator.ScaleAlerts {
		weightMap[alertDef.Name] = alertDef.Weight
	}

	scaEvtAggCfg := scaleAlertAggregator.Config{
		Logger:                 loggingFactory.NewNamedLogger("sokar.scaAlertAggr"),
		NoAlertScaleDamping:    cfg.ScaleAlertAggregator.NoAlertScaleDamping,
		UpScalingThreshold:     cfg.ScaleAlertAggregator.UpScaleThreshold,
		DownScalingThreshold:   cfg.ScaleAlertAggregator.DownScaleThreshold,
		EvaluationCycle:        cfg.ScaleAlertAggregator.EvaluationCycle,
		EvaluationPeriodFactor: cfg.ScaleAlertAggregator.EvaluationPeriodFactor,
		CleanupCycle:           cfg.ScaleAlertAggregator.CleanupCycle,
		WeightMap:              weightMap,
	}

	scaAlertAggr := scaEvtAggCfg.New(scaleAlertEmitters)
	api.Router.POST("/alert", scaAlertAggr.ScaleEvent)

	capaCfg := capacityPlanner.Config{
		Logger: loggingFactory.NewNamedLogger("sokar.capaPlanner"),
	}
	capaPlanner := capaCfg.New()

	sokarInst, err := setupSokar(scaAlertAggr, capaPlanner, scaler, api, logger)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed creating sokar.")
	}

	logger.Info().Msg("Connecting components and setting up sokar ... done")

	// Run all components
	sokarInst.Run()
	scaler.Run()
	scaAlertAggr.Run()
	api.Run()

	// Install signal handler for shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signalChan
		logger.Info().Msgf("Received %v. Shutting down...", s)

		// Stop all components
		api.Stop()
		scaAlertAggr.Stop()
		scaler.Stop()
		sokarInst.Stop()
	}()

	// Wait till completion
	api.Join()
	scaAlertAggr.Join()
	sokarInst.Join()

	logger.Info().Msg("Shutdown successfully completed")
	os.Exit(0)
}

func setupSokar(scaleEventEmitter sokarIF.ScaleEventEmitter, capacityPlanner sokarIF.CapacityPlanner, scaler sokarIF.Scaler, api api.API, logger zerolog.Logger) (*sokar.Sokar, error) {
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

func setupScaler(jobName string, min uint, max uint, nomadSrvAddr string, logF logging.LoggerFactory) (*scaler.Scaler, error) {

	// Set up the nomad connector
	nomadConfig := nomad.NewDefaultConfig(nomadSrvAddr)
	nomadConfig.Logger = logF.NewNamedLogger("sokar.nomad")
	nomad, err := nomadConfig.New()
	if err != nil {
		return nil, fmt.Errorf("Failed setting up nomad connector: %s.", err)
	}

	scaCfg := scaler.Config{
		JobName:  jobName,
		MinCount: min,
		MaxCount: max,
		Logger:   logF.NewNamedLogger("sokar.scaler"),
	}

	scaler, err := scaCfg.New(nomad)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up scaler: %s.", err)
	}

	return scaler, nil
}
