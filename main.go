package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/logging"
	"github.com/thomasobenaus/sokar/nomadConnector"
	"github.com/thomasobenaus/sokar/scaler"
	"github.com/thomasobenaus/sokar/sokar"
)

func main() {

	// parse commandline args and consume environment variables
	parsedArgs := parseArgs()
	if !parsedArgs.validateArgs() {
		os.Exit(1)
	}

	jobname := parsedArgs.JobName
	jobMinCount := parsedArgs.JobMinCount
	jobMaxCount := parsedArgs.JobMaxCount
	scaleBy := parsedArgs.ScaleBy
	localPort := 11000

	// set up logging
	lCfg := logging.Config{
		UseStructuredLogging:       parsedArgs.StructuredLogging,
		UseUnixTimestampForLogging: parsedArgs.UseUnixTimestampForLogging,
	}
	loggingFactory := lCfg.New()
	logger := loggingFactory.NewNamedLogger("sokar")

	logger.Info().Msg("Set up the scaler ...")
	scaler, err := setupScaler(jobname, jobMinCount, jobMaxCount, parsedArgs.NomadServerAddr, loggingFactory)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed setting up the scaler")
	}
	logger.Info().Msg("Set up the scaler ... done")

	// OneShot mode
	if parsedArgs.OneShot {
		scaResult := scaler.ScaleBy(scaleBy)
		logger.Info().Msgf("Scale %s: %s", scaResult.State, scaResult.StateDescription)
		os.Exit(0)
	}

	logger.Info().Msg("Connecting components and setting up sokar ...")
	sokarInst, err := setupSokar(scaler, logger)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed creating sokar.")
	}
	logger.Info().Msg("Connecting components and setting up sokar ... done")

	logger.Info().Msg("Registering http handlers ...")
	router := httprouter.New()
	router.POST(sokar.PathScaleBy, sokarInst.ScaleBy)

	logger.Info().Msg("Registering http handlers ... done")

	// Set up the web server
	logger.Info().Msgf("Start listening at %d.", localPort)
	if err := http.ListenAndServe(":"+strconv.Itoa(localPort), router); err != nil {
		logger.Fatal().Err(err).Msg("Failed serving.")
	}
}

func setupSokar(scaler sokar.Scaler, logger zerolog.Logger) (*sokar.Sokar, error) {
	cfg := sokar.Config{
		Logger: logger,
	}
	return cfg.New(scaler)
}

func setupScaler(jobName string, min uint, max uint, nomadSrvAddr string, logF logging.LoggerFactory) (*scaler.Scaler, error) {

	// Set up the nomad connector
	nomadConnectorConfig := nomadConnector.NewDefaultConfig(nomadSrvAddr)
	nomadConnectorConfig.Logger = logF.NewNamedLogger("sokar.nomad")
	nomadConnector, err := nomadConnectorConfig.New()
	if err != nil {
		return nil, fmt.Errorf("Failed setting up nomad connector: %s.", err)
	}

	scaCfg := scaler.Config{
		JobName:  jobName,
		MinCount: min,
		MaxCount: max,
		Logger:   logF.NewNamedLogger("sokar.scaler"),
	}

	scaler, err := scaCfg.New(nomadConnector)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up scaler: %s.", err)
	}

	return scaler, nil
}
