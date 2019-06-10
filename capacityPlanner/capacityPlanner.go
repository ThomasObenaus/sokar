package capacityPlanner

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

// CapaPlanningMode represents the mode the CapacityPlanner should use.
type CapaPlanningMode string

const (
	// CapaPlanningModeConstant specifies the mode where a constant offset is used to calculate the new planned scale.
	CapaPlanningModeConstant CapaPlanningMode = "const"
	// CapaPlanningModeLinear specifies the mode where the given scale is increased linearly based on the given scaleFactor.
	// Therefore the scaleFactor is directly used to scale the number of currentScale by multiplication.
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

	// constantMode if specified (not nil) the CapacityPlanner uses a constant offset to calculate the new planned scale. It is only allowed to
	// specify (not nil) one planning mode at the same time.
	constantMode *ConstantMode

	// linearMode if specified (not nil) the CapacityPlanner will increase the given scale linearly based on the current scaleFactor.
	// Therefore the scaleFactor is directly used to scale the number of currentScale by multiplication.
	// It is only allowed to specify (not nil) one planning mode at the same time.
	linearMode *LinearMode
}

// Config is the configuration for the Capacity Planner
type Config struct {
	Logger zerolog.Logger

	DownScaleCooldownPeriod time.Duration
	UpScaleCooldownPeriod   time.Duration

	ConstantMode *ConstantMode
	LinearMode   *LinearMode
}

// ConstantMode in this mode the CapacityPlanner uses a constant offset to calculate the new planned scale.
type ConstantMode struct {
	// Offset is the offset is just added/ substracted from the current scale to calculate the new planned scale.
	Offset uint
}

// LinearMode in this mode the CapacityPlanner will increase the given scale linearly based on the current scaleFactor.
// Therefore the scaleFactor is directly used to scale the number of currentScale by multiplication.
type LinearMode struct {
}

// NewDefaultConfig provides a config with good default values for the CapacityPlanner
func NewDefaultConfig() Config {
	return Config{
		DownScaleCooldownPeriod: time.Second * 80,
		UpScaleCooldownPeriod:   time.Second * 60,
		ConstantMode:            &ConstantMode{Offset: 1},
		LinearMode:              nil,
	}
}

// New creates a new instance of a CapacityPlanner using the given
// Scaler to send scaling events to.
func (cfg Config) New() (*CapacityPlanner, error) {

	if cfg.ConstantMode == nil && cfg.LinearMode == nil {
		return nil, fmt.Errorf("No planning mode specified")
	}

	if cfg.ConstantMode != nil && cfg.LinearMode != nil {
		return nil, fmt.Errorf("Multiple planning modes specified at the same time")
	}

	return &CapacityPlanner{
		logger:                  cfg.Logger,
		downScaleCooldownPeriod: cfg.DownScaleCooldownPeriod,
		upScaleCooldownPeriod:   cfg.UpScaleCooldownPeriod,
		constantMode:            cfg.ConstantMode,
		linearMode:              cfg.LinearMode,
	}, nil
}
