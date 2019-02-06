package main

import (
	"fmt"
	"os"

	"github.com/thomasobenaus/sokar/logging"
	"github.com/thomasobenaus/sokar/nomadConnector"
	"github.com/thomasobenaus/sokar/scaler"
)

func main() {

	// parse commandline args and consume environment variables
	parsedArgs := parseArgs()
	if !parsedArgs.validateArgs() {
		os.Exit(1)
	}

	jobname := parsedArgs.JobName
	scaleBy := parsedArgs.ScaleBy

	// set up logging
	lCfg := logging.Config{
		UseStructuredLogging:       parsedArgs.StructuredLogging,
		UseUnixTimestampForLogging: parsedArgs.UseUnixTimestampForLogging,
	}
	loggingFactory := lCfg.New()
	logger := loggingFactory.NewNamedLogger("sokar")

	logger.Info().Msg("Set up the scaler ...")
	scaler, err := setupScaler(jobname, parsedArgs.JobMinCount, parsedArgs.JobMaxCount, parsedArgs.NomadServerAddr, loggingFactory)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed setting up the scaler")
	}
	logger.Info().Msg("Set up the scaler ... done")

	err = scaler.ScaleBy(scaleBy)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to scale.")
	}
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
