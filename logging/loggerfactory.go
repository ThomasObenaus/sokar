package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// LoggerFactory is a factory that can be used to create named loggers using the same aligned configuration and namespace.
type LoggerFactory interface {
	NewNamedLogger(name string) zerolog.Logger
}

// New creates a new LoggerFactory which then can be used to create configured named loggers (log channels)
func New(structuredLogging, unixTimeStamp, disableColoredLogs bool) LoggerFactory {

	// default format for the timestamp
	zerolog.TimeFieldFormat = time.StampMilli //time.RFC3339

	if unixTimeStamp {
		// UNIX Time is faster and smaller than most timestamps
		// If you set zerolog.TimeFieldFormat to an empty string,
		// logs will write with UNIX time
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	return &loggerFactoryImpl{structuredLogging: structuredLogging, disableColoredLogs: disableColoredLogs}
}

type loggerFactoryImpl struct {
	structuredLogging  bool
	disableColoredLogs bool
}

// NewNamedLogger creates a new named and configured log-channel (logger)
func (lf *loggerFactoryImpl) NewNamedLogger(name string) zerolog.Logger {

	if lf.structuredLogging {
		return zerolog.New(os.Stdout).With().Timestamp().Str("logger", name).Logger()
	}

	return zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{
		NoColor: lf.disableColoredLogs, Out: os.Stderr,
		TimeFormat: zerolog.TimeFieldFormat,
	}).With().Timestamp().Str("logger", name).Logger()
}
