package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThomasObenaus/go-base/logging"
	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/alertmanager"
	"github.com/thomasobenaus/sokar/api"
	apipkg "github.com/thomasobenaus/sokar/api"
	"github.com/thomasobenaus/sokar/awsEc2"
	"github.com/thomasobenaus/sokar/capacityplanner"
	"github.com/thomasobenaus/sokar/config"
	"github.com/thomasobenaus/sokar/helper"
	"github.com/thomasobenaus/sokar/nomad"
	"github.com/thomasobenaus/sokar/nomadWorker"
	"github.com/thomasobenaus/sokar/scaleAlertAggregator"
	"github.com/thomasobenaus/sokar/scaler"
	"github.com/thomasobenaus/sokar/scaleschedule"
	"github.com/thomasobenaus/sokar/sokar"
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"

	"github.com/ThomasObenaus/go-base/health"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const endPointKey = "end-point"

func main() {
	// Print the build information as soon as possible to get at least some information on crashes
	buildinfo.Print(fmt.Printf)

	// read config
	cfg := helper.Must(cliAndConfig(os.Args)).(*config.Config)

	// set up logging
	loggingFactory := helper.Must(setupLogging(cfg)).(logging.LoggerFactory)

	logger := loggingFactory.NewNamedLogger("sokar")
	logger.Info().Msg("Connecting components and setting up sokar")

	logger.Info().Msg("1. Setup: API")
	api := apipkg.New(cfg.Port, apipkg.WithLogger(loggingFactory.NewNamedLogger("sokar.api")))

	logger.Info().Msg("2. Setup: ScaleSchedule")
	schedule := helper.Must(setupSchedule(cfg, logger)).(*scaleschedule.Schedule)

	logger.Info().Msg("3. Setup: ScaleAlertEmitters")
	scaleAlertEmitters := helper.Must(setupScaleAlertEmitters(api, loggingFactory)).([]scaleAlertAggregator.ScaleAlertEmitter)

	logger.Info().Msg("4. Setup: ScaleAlertAggregator")
	scaAlertAggr := setupScaleAlertAggregator(scaleAlertEmitters, cfg, loggingFactory)

	logger.Info().Msg("5. Setup: Scaling Target")
	scalingTarget := helper.Must(setupScalingTarget(cfg.Scaler, loggingFactory)).(scaler.ScalingTarget)
	logger.Info().Msgf("Scaling Target: %s", scalingTarget.String())

	logger.Info().Msg("6. Setup: Scaler")
	scaler := helper.Must(setupScaler(cfg.ScaleObject.Name, cfg.ScaleObject.MinCount, cfg.ScaleObject.MaxCount, cfg.Scaler.WatcherInterval, scalingTarget, loggingFactory, cfg.DryRunMode)).(*scaler.Scaler)

	logger.Info().Msg("7. Setup: CapacityPlanner")

	var mode capacityplanner.Option
	if cfg.CapacityPlanner.ConstantMode.Enable {
		mode = capacityplanner.UseConstantMode(cfg.CapacityPlanner.ConstantMode.Offset)
	} else if cfg.CapacityPlanner.LinearMode.Enable {
		mode = capacityplanner.UseLinearMode(float32(cfg.CapacityPlanner.LinearMode.ScaleFactorWeight))
	}

	capaPlanner := helper.Must(capacityplanner.New(
		capacityplanner.NewMetrics(),
		capacityplanner.WithLogger(loggingFactory.NewNamedLogger("sokar.capaPlanner")),
		capacityplanner.WithDownScaleCooldown(cfg.CapacityPlanner.DownScaleCooldownPeriod),
		capacityplanner.WithUpScaleCooldown(cfg.CapacityPlanner.UpScaleCooldownPeriod),
		capacityplanner.Schedule(schedule),
		mode,
	)).(*capacityplanner.CapacityPlanner)

	logger.Info().Msg("8. Setup: Sokar")
	sokarInst := helper.Must(setupSokar(scaAlertAggr, capaPlanner, scaler, schedule, api, logger, cfg.DryRunMode)).(*sokar.Sokar)

	// Setup health endpoint
	logger.Info().Msg("9. Setup: Health-Endpoint")
	loggerHealth := loggingFactory.NewNamedLogger("sokar.health")
	healthMonitor := helper.Must(health.NewMonitor(health.WithLogger(loggerHealth))).(*health.Monitor)
	healthMonitor.Start()
	if err := healthMonitor.Register(sokarInst); err != nil {
		logger.Fatal().Err(err).Msg("Failed to register health check for sokar")
	}
	api.Router.GET(sokar.PathHealth, apipkg.WrappedHandleFunc(healthMonitor.Health))
	logger.Info().Str(endPointKey, "health").Msgf("Health end-point set up at %s", sokar.PathHealth)

	// Register metrics handler
	api.Router.Handler("GET", sokar.PathMetrics, promhttp.Handler())
	logger.Info().Str(endPointKey, "metrics").Msgf("Metrics end-point set up at %s", sokar.PathMetrics)

	// Register build info end-point
	api.Router.GET(sokar.PathBuildInfo, apipkg.WrappedHandleFunc(buildinfo.BuildInfo))
	logger.Info().Str(endPointKey, "build info").Msgf("Build Info end-point set up at %s", sokar.PathBuildInfo)

	// Register config end-point
	cfgEndPoint := config.EndPoint{
		Config: *cfg,
		Logger: logger,
	}
	api.Router.GET(sokar.PathConfig, cfgEndPoint.ConfigEndpoint)
	logger.Info().Str(endPointKey, "config").Msgf("Config end-point set up at %s", sokar.PathConfig)

	// Define runnables and their execution order
	var orderedRunnables []Runnable
	orderedRunnables = append(orderedRunnables, sokarInst)
	orderedRunnables = append(orderedRunnables, scaler)
	orderedRunnables = append(orderedRunnables, scaAlertAggr)
	orderedRunnables = append(orderedRunnables, api)
	orderedRunnables = append(orderedRunnables, healthMonitor)

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
	cfg, err := config.New(args, "SK")
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

	loggingFactory := logging.New(cfg.Logging.Structured, cfg.Logging.UxTimestamp, cfg.Logging.NoColoredLogOutput)
	return loggingFactory, nil
}

func setupScaleAlertAggregator(scaleAlertEmitters []scaleAlertAggregator.ScaleAlertEmitter, cfg *config.Config, logF logging.LoggerFactory) *scaleAlertAggregator.ScaleAlertAggregator {
	weightMap := make(scaleAlertAggregator.ScaleAlertWeightMap)
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
	var scaleAlertEmitters []scaleAlertAggregator.ScaleAlertEmitter

	// Alertmanger Connector
	logger := logF.NewNamedLogger("sokar.alertmanager")
	amConnector := alertmanager.New(alertmanager.WithLogger(logger))
	api.Router.POST(sokar.PathAlertmanager, amConnector.HandleScaleAlerts)
	logger.Info().Msgf("Connector for alerts from prometheus/alertmanager setup successfully. Will listen for alerts on %s", sokar.PathAlertmanager)
	scaleAlertEmitters = append(scaleAlertEmitters, amConnector)

	return scaleAlertEmitters, nil
}

func setupSokar(scaleEventEmitter sokarIF.ScaleEventEmitter, capacityPlanner sokarIF.CapacityPlanner, scaler sokarIF.Scaler, schedule sokarIF.ScaleSchedule, api *api.API, logger zerolog.Logger, dryRunMode bool) (*sokar.Sokar, error) {
	cfg := sokar.Config{
		Logger:     logger,
		DryRunMode: dryRunMode,
	}
	sokarInst, err := cfg.New(scaleEventEmitter, capacityPlanner, scaler, schedule, sokar.NewMetrics())
	if err != nil {
		return nil, err
	}

	api.Router.PUT(sokar.PathScaleByPercentage, sokarInst.ScaleByPercentage)
	logger.Info().Str(endPointKey, "scale-by(p)").Msgf("ScaleBy end-point (percentage) set up at %s", sokar.PathScaleByPercentage)

	api.Router.PUT(sokar.PathScaleByValue, sokarInst.ScaleByValue)
	logger.Info().Str(endPointKey, "scale-by(v)").Msgf("ScaleBy end-point (value) set up at %s", sokar.PathScaleByValue)

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
		nomadWorker, err := nomadWorker.New(
			cfg.Nomad.ServerAddr,
			cfg.Nomad.DataCenterAWS.Profile,
			nomadWorker.WithLogger(logF.NewNamedLogger("sokar.nomadWorker")),
			nomadWorker.WithAwsRegion(cfg.Nomad.DataCenterAWS.Region),
			nomadWorker.TimeoutForInstanceTermination(cfg.Nomad.DataCenterAWS.InstanceTerminationTimeout),
		)
		if err != nil {
			return nil, fmt.Errorf("Failed setting up nomad worker connector: %s", err)
		}
		scalingTarget = nomadWorker
	} else if cfg.Mode == config.ScalerModeAwsEc2 {
		awsEc2, err := awsEc2.New(
			cfg.AwsEc2.ASGTagKey,
			cfg.AwsEc2.Profile,
			awsEc2.WithLogger(logF.NewNamedLogger("sokar.aws-ec2")),
			awsEc2.WithAwsRegion(cfg.AwsEc2.Region),
		)
		if err != nil {
			return nil, fmt.Errorf("Failed setting up aws-ec2 connector: %s", err)
		}
		scalingTarget = awsEc2
	} else {
		nomad, err := nomad.New(
			cfg.Nomad.ServerAddr,
			nomad.WithLogger(logF.NewNamedLogger("sokar.nomad")),
		)
		if err != nil {
			return nil, fmt.Errorf("Failed setting up nomad connector: %s", err)
		}

		scalingTarget = nomad
	}

	return scalingTarget, nil
}

// setupScaler creates and configures the Scaler. Internally nomad is used as scaling target.
func setupScaler(scalingObjName string, min uint, max uint, watcherInterval time.Duration, scalingTarget scaler.ScalingTarget, logF logging.LoggerFactory, dryRunMode bool) (*scaler.Scaler, error) {

	if logF == nil {
		return nil, fmt.Errorf("Logging factory is nil")
	}

	if scalingTarget == nil {
		return nil, fmt.Errorf("ScalingTarget is nil")
	}

	scalingObject := scaler.ScalingObject{Name: scalingObjName, MinCount: min, MaxCount: max}
	scaler, err := scaler.New(
		scalingTarget,
		scalingObject,
		scaler.NewMetrics(),
		scaler.WithLogger(logF.NewNamedLogger("sokar.scaler")),
		scaler.WatcherInterval(watcherInterval),
		scaler.DryRunMode(dryRunMode),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up scaler: %s", err)
	}

	return scaler, nil
}

// TODO: Add endpoint to provide schedule
func setupSchedule(cfg *config.Config, logger zerolog.Logger) (*scaleschedule.Schedule, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Config is nil")
	}

	scaleSchedule := scaleschedule.New()
	for _, entry := range cfg.CapacityPlanner.ScaleSchedule {
		for _, day := range entry.Days {
			minScale := uint(entry.MinScale)
			maxScale := uint(entry.MaxScale)
			if entry.MinScale < 0 {
				minScale = 0
			}

			if entry.MaxScale < 0 {
				maxScale = helper.MaxUint
			}
			if err := scaleSchedule.Insert(day, entry.StartTime, entry.EndTime, minScale, maxScale); err != nil {
				logger.Warn().Msgf("Entry '%s' was not added to scale schedule for %s: %s", entry, day, err.Error())
			} else {
				logger.Debug().Msgf("Entry to scale schedule added: On %s at %s to %s -> [%d,%d]", day, entry.StartTime, entry.EndTime, minScale, maxScale)
			}
		}

	}

	return &scaleSchedule, nil
}
