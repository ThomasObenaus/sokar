package capacityPlanner

import (
	"github.com/rs/zerolog"
)

// CapacityPlanner is a object that plans new resources based on current needs and directly
type CapacityPlanner struct {
	logger zerolog.Logger
}

// Config is the configuration for the Capacity Planner
type Config struct {
	Logger zerolog.Logger
}

// New creates a new instance of a CapacityPlanner using the given
// Scaler to send scaling events to.
func (cfg Config) New() *CapacityPlanner {
	return &CapacityPlanner{
		logger: cfg.Logger,
	}
}
