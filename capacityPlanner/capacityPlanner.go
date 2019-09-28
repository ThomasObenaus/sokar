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

	// constantMode if specified (not nil) the CapacityPlanner uses a constant offset to calculate the new planned scale. It is only allowed to
	// specify (not nil) one planning mode at the same time.
	constantMode *ConstantMode

	// linearMode if specified (not nil) the CapacityPlanner will increase the given scale linearly based on the current scaleFactor.
	// Therefore the scaleFactor is directly used to scale the number of currentScale by multiplication.
	// It is only allowed to specify (not nil) one planning mode at the same time.
	linearMode *LinearMode
}

// ConstantMode in this mode the CapacityPlanner uses a constant offset to calculate the new planned scale.
type ConstantMode struct {
	// Offset is the offset is just added/ substracted from the current scale to calculate the new planned scale.
	Offset uint
}

// LinearMode in this mode the CapacityPlanner will increase the given scale linearly based on the current scaleFactor.
// Therefore the scaleFactor is directly used to scale the number of currentScale by multiplication.
type LinearMode struct {
	ScaleFactorWeight float32
}

// Option represents an option for the CapacityPlanner
type Option func(cp *CapacityPlanner)

// WithLogger adds a configured Logger to the CapacityPlanner
func WithLogger(logger zerolog.Logger) Option {
	return func(cp *CapacityPlanner) {
		cp.logger = logger
	}
}

// UseConstantMode switches the CapacityPlanner to constant mode
func UseConstantMode(offset uint) Option {
	return func(cp *CapacityPlanner) {
		cp.constantMode = &ConstantMode{Offset: offset}
		cp.linearMode = nil
	}
}

// UseLinearMode switches the CapacityPlanner to linear mode
func UseLinearMode(scaleFactorWeight float32) Option {
	return func(cp *CapacityPlanner) {
		cp.linearMode = &LinearMode{ScaleFactorWeight: scaleFactorWeight}
		cp.constantMode = nil
	}
}

// WithDownScaleCooldown specifies the cooldown of the CapacityPlanner after downscale action
func WithDownScaleCooldown(cooldown time.Duration) Option {
	return func(cp *CapacityPlanner) {
		cp.downScaleCooldownPeriod = cooldown
	}
}

// WithUpScaleCooldown specifies the cooldown of the CapacityPlanner after upscale action
func WithUpScaleCooldown(cooldown time.Duration) Option {
	return func(cp *CapacityPlanner) {
		cp.upScaleCooldownPeriod = cooldown
	}
}

// New creates a new instance of a CapacityPlanner using the given
// Scaler to send scaling events to.
func New(options ...Option) (*CapacityPlanner, error) {

	capacityPlanner := CapacityPlanner{
		downScaleCooldownPeriod: time.Second * 80,
		upScaleCooldownPeriod:   time.Second * 60,
		constantMode:            &ConstantMode{Offset: 1},
		linearMode:              nil,
	}

	// apply the options
	for _, opt := range options {
		opt(&capacityPlanner)
	}

	if err := validate(capacityPlanner); err != nil {
		return nil, err
	}

	return &capacityPlanner, nil
}

func validate(cp CapacityPlanner) error {
	if cp.constantMode == nil && cp.linearMode == nil {
		return fmt.Errorf("No planning mode specified")
	}

	if cp.constantMode != nil && cp.linearMode != nil {
		return fmt.Errorf("Multiple planning modes specified at the same time")
	}

	if cp.linearMode != nil && cp.linearMode.ScaleFactorWeight <= 0 {
		return fmt.Errorf("The given value for the ScaleFactorWeight '%f' is not allowed. Only values > 0 are valid", cp.linearMode.ScaleFactorWeight)
	}

	if cp.constantMode != nil && cp.constantMode.Offset == 0 {
		return fmt.Errorf("The given value for the offset '%d' is not allowed. Only values > 0 are valid", cp.constantMode.Offset)
	}

	return nil
}
