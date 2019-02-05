package main

import (
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

	jobname := "fail-service"

	// set up logging
	lCfg := logging.Config{
		UseStructuredLogging:       parsedArgs.StructuredLogging,
		UseUnixTimestampForLogging: parsedArgs.UseUnixTimestampForLogging,
	}
	loggingFactory := lCfg.New()
	logger := loggingFactory.NewNamedLogger("sokar")

	logger.Info().Msg("Set up the scaler ...")

	// Set up the nomad connector
	nomadConnectorConfig := nomadConnector.NewDefaultConfig(parsedArgs.NomadServerAddr)
	nomadConnectorConfig.Logger = loggingFactory.NewNamedLogger("sokar.nomad")
	nomadConnector, err := nomadConnectorConfig.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed setting up nomad connector")
	}

	scaCfg := scaler.Config{
		JobName:  jobname,
		MinCount: 1,
		MaxCount: 10,
		Logger:   loggingFactory.NewNamedLogger("sokar.scaler"),
	}

	scaler, err := scaCfg.New(nomadConnector)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed setting up scaler")
	}

	logger.Info().Msg("Set up the scaler ... done")

	err = scaler.ScaleBy(-5)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to scale.")
	}
}
