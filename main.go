package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/alertmanager"
	"github.com/thomasobenaus/sokar/api"
	"github.com/thomasobenaus/sokar/awsEc2"
	"github.com/thomasobenaus/sokar/capacityPlanner"
	"github.com/thomasobenaus/sokar/config"
	"github.com/thomasobenaus/sokar/helper"
	"github.com/thomasobenaus/sokar/logging"
	"github.com/thomasobenaus/sokar/nomad"
	"github.com/thomasobenaus/sokar/nomadWorker"
	"github.com/thomasobenaus/sokar/scaleAlertAggregator"
	"github.com/thomasobenaus/sokar/scaler"
	"github.com/thomasobenaus/sokar/sokar"
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var version string
var buildTime string
var revision string
var branch string

func main() {

	// read config
	cfg := helper.Must(cliAndConfig(os.Args)).(*config.Config)

	buildInfo := BuildInfo{
		Version:   version,
		BuildTime: buildTime,
		Revision:  revision,
		Branch:    branch,
	}
	buildInfo.Print(fmt.Printf)

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

	logger.Info().Msg("4. Setup: Scaling Target")
	scalingTarget := helper.Must(setupScalingTarget(cfg.Scaler, loggingFactory)).(scaler.ScalingTarget)
	logger.Info().Msgf("Scaling Target: %s", scalingTarget.String())

	logger.Info().Msg("5. Setup: Scaler")
	scaler := helper.Must(setupScaler(cfg.ScaleObject.Name, cfg.ScaleObject.MinCount, cfg.ScaleObject.MaxCount, cfg.Scaler.WatcherInterval, scalingTarget, loggingFactory)).(*scaler.Scaler)

	logger.Info().Msg("6. Setup: CapacityPlanner")

	var constantMode *capacityPlanner.ConstantMode
	var linearMode *capacityPlanner.LinearMode
	if cfg.CapacityPlanner.ConstantMode.Enable {
		constantMode = &capacityPlanner.ConstantMode{Offset: cfg.CapacityPlanner.ConstantMode.Offset}
	} else if cfg.CapacityPlanner.LinearMode.Enable {
		linearMode = &capacityPlanner.LinearMode{ScaleFactorWeight: float32(cfg.CapacityPlanner.LinearMode.ScaleFactorWeight)}
	}

	capaCfg := capacityPlanner.Config{
		Logger:                  loggingFactory.NewNamedLogger("sokar.capaPlanner"),
		DownScaleCooldownPeriod: cfg.CapacityPlanner.DownScaleCooldownPeriod,
		UpScaleCooldownPeriod:   cfg.CapacityPlanner.UpScaleCooldownPeriod,
		ConstantMode:            constantMode,
		LinearMode:              linearMode,
	}
	capaPlanner := helper.Must(capaCfg.New()).(*capacityPlanner.CapacityPlanner)

	logger.Info().Msg("7. Setup: Sokar")
	sokarInst := helper.Must(setupSokar(scaAlertAggr, capaPlanner, scaler, api, logger, cfg.DryRunMode)).(*sokar.Sokar)

	// Register metrics handler
	api.Router.Handler("GET", sokar.PathMetrics, promhttp.Handler())
	logger.Info().Str("end-point", "metrics").Msgf("Metrics end-point set up at %s", sokar.PathMetrics)

	// Register build info end-point
	api.Router.GET(sokar.PathBuildInfo, buildInfo.BuildInfo)
	logger.Info().Str("end-point", "build info").Msgf("Build Info end-point set up at %s", sokar.PathBuildInfo)

	// Register config end-point
	cfgEndPoint := config.EndPoint{
		Config: *cfg,
		Logger: logger,
	}
	api.Router.GET(sokar.PathConfig, cfgEndPoint.ConfigEndpoint)
	logger.Info().Str("end-point", "config").Msgf("Config end-point set up at %s", sokar.PathConfig)

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
	// and read config
	cfg := config.NewDefaultConfig()
	err := cfg.ReadConfig(args)
	if err != nil {
		return nil, err
	}

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
		NoColoredLogOutput:         cfg.Logging.NoColoredLogOutput,
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
	logger := logF.NewNamedLogger("sokar.alertmanager")
	amCfg := alertmanager.Config{
		Logger: logger,
	}
	amConnector := amCfg.New()
	api.Router.POST(sokar.PathAlertmanager, amConnector.HandleScaleAlerts)
	logger.Info().Msgf("Connector for alerts from prometheus/alertmanager setup successfully. Will listen for alerts on %s", sokar.PathAlertmanager)

	var scaleAlertEmitters []scaleAlertAggregator.ScaleAlertEmitter
	scaleAlertEmitters = append(scaleAlertEmitters, amConnector)

	return scaleAlertEmitters, nil
}

func setupSokar(scaleEventEmitter sokarIF.ScaleEventEmitter, capacityPlanner sokarIF.CapacityPlanner, scaler sokarIF.Scaler, api *api.API, logger zerolog.Logger, dryRunMode bool) (*sokar.Sokar, error) {
	cfg := sokar.Config{
		Logger:     logger,
		DryRunMode: dryRunMode,
	}
	sokarInst, err := cfg.New(scaleEventEmitter, capacityPlanner, scaler, sokar.NewMetrics())
	if err != nil {
		return nil, err
	}

	api.Router.GET(sokar.PathHealth, sokarInst.Health)
	logger.Info().Str("end-point", "health").Msgf("Health end-point set up at %s", sokar.PathHealth)

	api.Router.PUT(sokar.PathScaleByPercentage, sokarInst.ScaleByPercentage)
	logger.Info().Str("end-point", "scale-by(p)").Msgf("ScaleBy end-point (percentage) set up at %s", sokar.PathScaleByPercentage)

	api.Router.PUT(sokar.PathScaleByValue, sokarInst.ScaleByValue)
	logger.Info().Str("end-point", "scale-by(v)").Msgf("ScaleBy end-point (value) set up at %s", sokar.PathScaleByValue)

	if cfg.DryRunMode {
		logger.Info().Msg("Dry-Run-Mode: Sokar will plan the scale actions but won't execute them. This applies only for auto scaling events.")
	}

	return sokarInst, err
}

func setupScalingTarget(cfg config.Scaler, logF logging.LoggerFactory) (scaler.ScalingTarget, error) {
	if logF == nil {
		return nil, fmt.Errorf("Logging factory is nil")
	}

	var scalingTarget scaler.ScalingTarget

	if cfg.Mode == config.ScalerModeNomadDataCenter {
		cfg := nomadWorker.Config{NomadServerAddress: cfg.Nomad.ServerAddr, Logger: logF.NewNamedLogger("sokar.nomadWorker"), AWSRegion: cfg.Nomad.DataCenterAWS.Region, AWSProfile: cfg.Nomad.DataCenterAWS.Profile}
		nomadWorker, err := cfg.New()
		if err != nil {
			return nil, fmt.Errorf("Failed setting up nomad worker connector: %s", err)
		}
		scalingTarget = nomadWorker
	} else if cfg.Mode == config.ScalerModeAwsEc2 {
		cfg := awsEc2.Config{Logger: logF.NewNamedLogger("sokar.aws-ec2"), AWSRegion: cfg.AwsEc2.Region, AWSProfile: cfg.AwsEc2.Profile, ASGTagKey: cfg.AwsEc2.ASGTagKey}
		awsEc2, err := cfg.New()
		if err != nil {
			return nil, fmt.Errorf("Failed setting up aws-ec2 connector: %s", err)
		}
		scalingTarget = awsEc2
	} else {
		nomadConfig := nomad.NewDefaultConfig(cfg.Nomad.ServerAddr)
		nomadConfig.Logger = logF.NewNamedLogger("sokar.nomad")
		nomad, err := nomadConfig.New()
		if err != nil {
			return nil, fmt.Errorf("Failed setting up nomad connector: %s", err)
		}

		scalingTarget = nomad
	}

	return scalingTarget, nil
}

// setupScaler creates and configures the Scaler. Internally nomad is used as scaling target.
func setupScaler(scalingObjName string, min uint, max uint, watcherInterval time.Duration, scalingTarget scaler.ScalingTarget, logF logging.LoggerFactory) (*scaler.Scaler, error) {

	if logF == nil {
		return nil, fmt.Errorf("Logging factory is nil")
	}

	if scalingTarget == nil {
		return nil, fmt.Errorf("ScalingTarget is nil")
	}

	scaCfg := scaler.Config{
		Name:            scalingObjName,
		MinCount:        min,
		MaxCount:        max,
		Logger:          logF.NewNamedLogger("sokar.scaler"),
		WatcherInterval: watcherInterval,
	}

	scaler, err := scaCfg.New(scalingTarget, scaler.NewMetrics())
	if err != nil {
		return nil, fmt.Errorf("Failed setting up scaler: %s", err)
	}

	return scaler, nil
}
