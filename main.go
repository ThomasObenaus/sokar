package main

func main() {

	// parse commandline args and consume environment variables
	parsedArgs := parseArgs()

	lCfg := LoggingCfg{
		LoggerName:                 "sokar",
		UseStructuredLogging:       parsedArgs.StructuredLogging,
		UseUnixTimestampForLogging: parsedArgs.UseUnixTimestampForLogging,
	}
	log := NewLogger(lCfg)

	log.Info().Float64("duration", 29.343).Str("region", "ED01").Msg("hello world")

}
