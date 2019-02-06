package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerFactory interface {
	NewNamedLogger(name string) zerolog.Logger
}

type Config struct {
	UseStructuredLogging       bool
	UseUnixTimestampForLogging bool
}

func (cfg Config) New() LoggerFactory {
	if cfg.UseUnixTimestampForLogging {
		// UNIX Time is faster and smaller than most timestamps
		// If you set zerolog.TimeFieldFormat to an empty string,
		// logs will write with UNIX time
		zerolog.TimeFieldFormat = ""
	} else {
		zerolog.TimeFieldFormat = time.StampMilli //time.RFC3339
	}

	return &loggerFactoryImpl{UseStructuredLogging: cfg.UseStructuredLogging}
}

type loggerFactoryImpl struct {
	UseStructuredLogging bool
}

func (lf *loggerFactoryImpl) NewNamedLogger(name string) zerolog.Logger {

	var logger zerolog.Logger
	if lf.UseStructuredLogging {
		logger = zerolog.New(os.Stdout).With().Timestamp().Str("logger", name).Logger()
	} else {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Str("logger", name).Logger()
	}

	return logger
}
