package scaler

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Scaler is a component responsible for scaling a job
type Scaler struct {
	logger        zerolog.Logger
	scalingTarget ScalingTarget
	job           jobConfig

	// jobWatcherCycle the cycle the Scaler will check if
	// the job count still matches the desired state.
	jobWatcherCycle time.Duration

	desiredCount uint

	numOpenScalingTickets uint
	maxOpenScalingTickets uint
	scaleTicketChan       chan ScalingTicket

	// channel used to signal teardown/ stop
	stopChan chan struct{}

	wg sync.WaitGroup
}

// Config is the configuration for the Scaler
type Config struct {
	JobName  string
	MinCount uint
	MaxCount uint
	Logger   zerolog.Logger
}

// New creates a new instance of a scaler using the given
// ScalingTarget to send scaling events to.
func (cfg Config) New(scalingTarget ScalingTarget) (*Scaler, error) {
	if scalingTarget == nil {
		return nil, fmt.Errorf("Given ScalingTarget is nil")
	}

	return &Scaler{
		logger:          cfg.Logger,
		scalingTarget:   scalingTarget,
		jobWatcherCycle: time.Second * 5,
		job: jobConfig{
			jobName:  cfg.JobName,
			minCount: cfg.MinCount,
			maxCount: cfg.MaxCount,
		},
		stopChan:              make(chan struct{}, 1),
		numOpenScalingTickets: 0,
		maxOpenScalingTickets: 0,
		scaleTicketChan:       make(chan ScalingTicket, 1),
	}, nil
}

// GetCount returns the number of currently deployed instances
func (s *Scaler) GetCount() (uint, error) {
	return s.scalingTarget.GetJobCount(s.job.jobName)
}

// ScaleTo will scale the job to the desired count.
func (s *Scaler) ScaleTo(desiredCount uint) error {
	s.logger.Info().Msgf("Scale to %d requested.", desiredCount)
	return s.openScalingTicket(desiredCount)
}

// Run starts/ runs the scaler
func (s *Scaler) Run() {
	// handler that processes incoming scaling tickets
	go s.scaleTicketProcessor(s.scaleTicketChan)
	// handler that checks periodically if the desired count is still valid
	go s.jobWatcher(s.jobWatcherCycle)
}

// Stop tears down scaler
func (s *Scaler) Stop() {
	s.logger.Info().Msg("Teardown requested")

	close(s.scaleTicketChan)
	close(s.stopChan)
}

// Join blocks/ waits until scaler has been stopped
func (s *Scaler) Join() {
	s.wg.Wait()
}
