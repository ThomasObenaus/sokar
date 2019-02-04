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
	lCfg := logging.LoggingCfg{
		LoggerName:                 "sokar",
		UseStructuredLogging:       parsedArgs.StructuredLogging,
		UseUnixTimestampForLogging: parsedArgs.UseUnixTimestampForLogging,
	}
	log := lCfg.New()

	nomadConnectorConfig := nomadConnector.Config{
		JobName:            "fail-service",
		NomadServerAddress: parsedArgs.NomadServerAddr,
	}
	nomadConnector, err := nomadConnectorConfig.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed setting up nomad connector")
	}

	nomadConnector.ScaleBy(2)

}
