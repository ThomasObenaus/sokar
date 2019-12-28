package alertscheduler

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	saa "github.com/thomasobenaus/sokar/scaleAlertAggregator"
)

// AlertScheduler is a component that emits ScaleAlerts based on a given schedule
type AlertScheduler struct {
	logger zerolog.Logger

	schedule AlertSchedule

	// handleFuncs is a list of registered handlers for received ScaleAlerts
	handleFuncs []saa.ScaleAlertHandleFunc

	// channel used to signal teardown/ stop
	stopChan        chan struct{}
	wg              sync.WaitGroup
	evaluationCycle time.Duration
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
func New(schedule AlertSchedule, options ...Option) *AlertScheduler {

	alertScheduler := &AlertScheduler{
		handleFuncs:     make([]saa.ScaleAlertHandleFunc, 0),
		schedule:        schedule,
		evaluationCycle: time.Second * 30,
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

// Run starts AlertScheduler
func (as *AlertScheduler) Run() {

	evaluationTicker := time.NewTicker(as.evaluationCycle)

	// main loop
	go func() {
		as.logger.Info().Msg("Main process loop started")

	loop:
		for {
			select {

			case <-as.stopChan:
				evaluationTicker.Stop()
				close(as.stopChan)
				break loop

			case <-evaluationTicker.C:
				as.logger.Info().Msg("####################")
			}
		}
		as.logger.Info().Msg("Main process loop left")
	}()
}

// Join blocks/ waits until AlertScheduler has been stopped
func (as *AlertScheduler) Join() {
	as.wg.Wait()
}

// Stop tears down AlertScheduler
func (as *AlertScheduler) Stop() error {
	as.logger.Info().Msg("Teardown requested")
	close(as.stopChan)
	return nil
}

// GetName returns the name of this component
func (as *AlertScheduler) GetName() string {
	return "AlertScheduler"
}
