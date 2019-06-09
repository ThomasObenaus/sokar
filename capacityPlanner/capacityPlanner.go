package capacityPlanner

import (
	"time"

	"github.com/rs/zerolog"
)

// CapaPlanningMode represents the mode the CapacityPlanner should use
type CapaPlanningMode string

const (
	// CapaPlanningModeConstant specifies the mode where a constant offset is used to calculate the new planned scale
	CapaPlanningModeConstant CapaPlanningMode = "const"
	// CapaPlanningModeLinear specifies the mode where the given scale is increased linearly based on the given scaleFactor. Therefore the scaleFactor is directly used to scale the number of currentScale by multiplication.
	CapaPlanningModeLinear CapaPlanningMode = "linear"
)

// CapacityPlanner is a object that plans new resources based on current needs and directly
type CapacityPlanner struct {
	logger                  zerolog.Logger
	downScaleCooldownPeriod time.Duration
	upScaleCooldownPeriod   time.Duration

	// mode specifies the mode that shall be used to calculate the new planned scale
	mode CapaPlanningMode

	// offsetConstantMode is the offset that is used in CapaPlanningModeConstant.
	// There this offset is just added/ substracted from the current scale to calculate the new planned scale.
	offsetConstantMode uint
}

// Config is the configuration for the Capacity Planner
type Config struct {
	Logger zerolog.Logger

	DownScaleCooldownPeriod time.Duration
	UpScaleCooldownPeriod   time.Duration

	Mode               CapaPlanningMode
	OffsetConstantMode uint
}

// NewDefaultConfig provides a config with good default values for the CapacityPlanner
func NewDefaultConfig() Config {
	return Config{
		DownScaleCooldownPeriod: time.Second * 80,
		UpScaleCooldownPeriod:   time.Second * 60,
		OffsetConstantMode:      1,
		Mode:                    CapaPlanningModeConstant,
	}
}

// New creates a new instance of a CapacityPlanner using the given
// Scaler to send scaling events to.
func (cfg Config) New() *CapacityPlanner {
	return &CapacityPlanner{
		logger:                  cfg.Logger,
		downScaleCooldownPeriod: cfg.DownScaleCooldownPeriod,
		upScaleCooldownPeriod:   cfg.UpScaleCooldownPeriod,
		offsetConstantMode:      cfg.OffsetConstantMode,
		mode:                    cfg.Mode,
	}
}
