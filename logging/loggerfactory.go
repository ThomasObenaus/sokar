package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LoggerFactory is a factory that can be used to create named loggers using the same aligned configuration and namespace.
type LoggerFactory interface {
	NewNamedLogger(name string) zerolog.Logger
}

// Config is a struct keeping the main configuration parameters for the LoggerFactory
type Config struct {
	UseStructuredLogging       bool
	UseUnixTimestampForLogging bool
	NoColoredLogOutput         bool
}

// New creates a new LoggerFactory
func (cfg Config) New() LoggerFactory {
	if cfg.UseUnixTimestampForLogging {
		// UNIX Time is faster and smaller than most timestamps
		// If you set zerolog.TimeFieldFormat to an empty string,
		// logs will write with UNIX time
		zerolog.TimeFieldFormat = ""
	} else {
		zerolog.TimeFieldFormat = time.StampMilli //time.RFC3339
	}

	return &loggerFactoryImpl{UseStructuredLogging: cfg.UseStructuredLogging, NoColoredLogOutput: cfg.NoColoredLogOutput}
}

type loggerFactoryImpl struct {
	UseStructuredLogging bool
	NoColoredLogOutput   bool
}

func (lf *loggerFactoryImpl) NewNamedLogger(name string) zerolog.Logger {

	var logger zerolog.Logger
	if lf.UseStructuredLogging {
		logger = zerolog.New(os.Stdout).With().Timestamp().Str("logger", name).Logger()
	} else {
		logger = log.Output(zerolog.ConsoleWriter{NoColor: lf.NoColoredLogOutput, Out: os.Stderr}).With().Timestamp().Str("logger", name).Logger()
	}

	return logger
}
