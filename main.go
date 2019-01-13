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
	nomadConnector := nomadConnectorConfig.New()

	nomadConnector.ScaleBy(2)

	log.Info().Float64("duration", 29.343).Str("region", "ED01").Msg("hello world")

}
