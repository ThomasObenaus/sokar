package scaler

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/helper"
	sokar "github.com/thomasobenaus/sokar/sokar/iface"
)

// Scaler is a component responsible for scaling a job
type Scaler struct {
	logger        zerolog.Logger
	scalingTarget ScalingTarget
	job           jobConfig

	// jobWatcherCycle the cycle the Scaler will check if
	// the job count still matches the desired state.
	jobWatcherCycle time.Duration

	scalingTicket *ScalingTicket

	// lock to sync the multi-thread
	// access to the ScalingTicket
	lock sync.RWMutex

	// channel used to signal teardown/ stop
	stopChan chan struct{}
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
		stopChan:      make(chan struct{}, 1),
		scalingTicket: nil,
	}, nil
}

// ScaleTo scales the job to the given count
func (s *Scaler) ScaleTo(desiredCount uint) sokar.ScaleResult {
	if r, ok := trueIfNil(s); ok {
		return r
	}

	jobName := s.job.jobName
	currentCount, err := s.scalingTarget.GetJobCount(jobName)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error obtaining job count: %s.", err.Error()),
		}
	}

	return s.scale(desiredCount, currentCount)
}

// ScaleBy Scales the target component by the given amount of instances
func (s *Scaler) ScaleBy(amount int) sokar.ScaleResult {
	if r, ok := trueIfNil(s); ok {
		return r
	}

	jobName := s.job.jobName
	count, err := s.scalingTarget.GetJobCount(jobName)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error obtaining job count: %s.", err.Error()),
		}
	}

	desiredCount := helper.IncUint(count, amount)

	return s.scale(desiredCount, count)
}
