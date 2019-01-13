package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LoggingCfg holds the configuration for logging
type LoggingCfg struct {
	LoggerName                 string
	UseStructuredLogging       bool
	UseUnixTimestampForLogging bool
}

// NewLogger creates a new logger object with a given log-channel name
func NewLogger(loggingCfg LoggingCfg) zerolog.Logger {

	if loggingCfg.UseUnixTimestampForLogging {
		// UNIX Time is faster and smaller than most timestamps
		// If you set zerolog.TimeFieldFormat to an empty string,
		// logs will write with UNIX time
		zerolog.TimeFieldFormat = ""
	}

	var logger zerolog.Logger
	if loggingCfg.UseStructuredLogging {
		logger = zerolog.New(os.Stdout).With().Timestamp().Str("logger", loggingCfg.LoggerName).Logger()
	} else {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Str("logger", loggingCfg.LoggerName).Logger()
	}

	return logger
}
