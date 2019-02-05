package scaler

import (
	"fmt"

	"github.com/rs/zerolog"
)

// Scaler specifies the interface for a component that can scale a certain job
type Scaler interface {
	// ScaleBy Scales the target component by the given amount of instances
	ScaleBy(amount int) error
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
func (cfg Config) New(scalingTarget ScalingTarget) (Scaler, error) {
	if scalingTarget == nil {
		return nil, fmt.Errorf("Given ScalingTarget is nil")
	}

	return &scalerImpl{
		logger:        cfg.Logger,
		scalingTarget: scalingTarget,
		job: jobConfig{
			jobName:  cfg.JobName,
			minCount: cfg.MinCount,
			maxCount: cfg.MaxCount,
		},
	}, nil
}
