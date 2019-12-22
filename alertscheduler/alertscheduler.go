package alertscheduler

import (
	"github.com/rs/zerolog"
	saa "github.com/thomasobenaus/sokar/scaleAlertAggregator"
)

// AlertScheduler is a component that emits ScaleAlerts based on a given schedule
type AlertScheduler struct {
	logger zerolog.Logger

	// handleFuncs is a list of registered handlers for received ScaleAlerts
	handleFuncs []saa.ScaleAlertHandleFunc
}

// Option represents an option for the alertscheduler
type Option func(as *AlertScheduler)

// WithLogger adds a configured Logger to the alertscheduler
func WithLogger(logger zerolog.Logger) Option {
	return func(as *AlertScheduler) {
		as.logger = logger
	}
}

// New creates a new instance of a AlertScheduler
func New(options ...Option) *AlertScheduler {

	alertScheduler := &AlertScheduler{
		handleFuncs: make([]saa.ScaleAlertHandleFunc, 0),
	}
	// apply the options
	for _, opt := range options {
		opt(alertScheduler)
	}

	return alertScheduler
}

// Register is used to register a handler/ listener for scaleAlertEvents (ScaleAlertPacket)
func (as *AlertScheduler) Register(handleFunc saa.ScaleAlertHandleFunc) {
	as.handleFuncs = append(as.handleFuncs, handleFunc)
}
