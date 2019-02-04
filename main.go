package main

import (
	"os"

	"github.com/thomasobenaus/sokar/logging"
	"github.com/thomasobenaus/sokar/nomadConnector"
)

func main() {

	// parse commandline args and consume environment variables
	parsedArgs := parseArgs()
	if !parsedArgs.validateArgs() {
		os.Exit(1)
	}

	// set up logging
	lcfg := logging.Config{
		UseStructuredLogging:       parsedArgs.StructuredLogging,
		UseUnixTimestampForLogging: parsedArgs.UseUnixTimestampForLogging,
	}
	loggingFactory := lcfg.New()
	logger := loggingFactory.NewNamedLogger("sokar")

	// Set up the nomad connector
	nomadConnectorConfig := nomadConnector.Config{
		JobName:            "fail-service",
		NomadServerAddress: parsedArgs.NomadServerAddr,
		Logger:             loggingFactory.NewNamedLogger("sokar.nomad"),
	}

	nomadConnector, err := nomadConnectorConfig.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed setting up nomad connector")
	}

	nomadConnector.ScaleBy(2)

}
