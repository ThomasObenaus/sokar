package capacityPlanner

import (
	"time"

	"github.com/rs/zerolog"
)

// CapacityPlanner is a object that plans new resources based on current needs and directly
type CapacityPlanner struct {
	logger                  zerolog.Logger
	downScaleCooldownPeriod time.Duration
	upScaleCooldownPeriod   time.Duration
}

// Config is the configuration for the Capacity Planner
type Config struct {
	Logger zerolog.Logger

	DownScaleCooldownPeriod time.Duration
	UpScaleCooldownPeriod   time.Duration
}

// NewDefaultConfig provides a config with good default values for the CapacityPlanner
func NewDefaultConfig() Config {
	return Config{
		DownScaleCooldownPeriod: time.Second * 80,
		UpScaleCooldownPeriod:   time.Second * 60,
	}
}

// New creates a new instance of a CapacityPlanner using the given
// Scaler to send scaling events to.
func (cfg Config) New() *CapacityPlanner {
	return &CapacityPlanner{
		logger:                  cfg.Logger,
		downScaleCooldownPeriod: cfg.DownScaleCooldownPeriod,
		upScaleCooldownPeriod:   cfg.UpScaleCooldownPeriod,
	}
}
