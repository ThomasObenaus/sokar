package main

import (
	"github.com/thomasobenaus/sokar/nomadConnector"
)

func main() {

	// parse commandline args and consume environment variables
	parsedArgs := parseArgs()

	// set up logging
	lCfg := LoggingCfg{
		LoggerName:                 "sokar",
		UseStructuredLogging:       parsedArgs.StructuredLogging,
		UseUnixTimestampForLogging: parsedArgs.UseUnixTimestampForLogging,
	}
	log := lCfg.New()

	nomadConnectorConfig := nomadConnector.Config{
		JobName: "ping-service",
	}
	nomadConnector := nomadConnectorConfig.New(log)

	nomadConnector.ScaleBy(2)

}
