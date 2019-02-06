package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

	if parsedArgs.OneShot {
		err = scaler.ScaleBy(scaleBy)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scale.")
		}
	} else {

		logger.Info().Msg("Connecting components and setting up sokar ...")
		if err := setupSokar(scaler, logger); err != nil {
			logger.Fatal().Err(err).Msg("Failed creating sokar.")
		}
		logger.Info().Msg("Connecting components and setting up sokar ... done")

		// Set up the web server
		logger.Info().Msgf("Start listening at %d.", localPort)
		if err := http.ListenAndServe(":"+strconv.Itoa(localPort), nil); err != nil {
			logger.Fatal().Err(err).Msg("Failed serving.")
		}
	}
}

func setupSokar(scaler sokar.Scaler, logger zerolog.Logger) error {
	cfg := sokar.Config{
		Logger: logger,
	}
	sokar, err := cfg.New(scaler)

	// Register the handlers provided by sokar
	http.HandleFunc("/scaler", sokar.HandleScaler)

	return err
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
