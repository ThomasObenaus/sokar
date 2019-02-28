package scaleAlertAggregator

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/helper"
	"github.com/thomasobenaus/sokar/sokar"
)

// ScaleAlertAggregator is a component that is responsible to gather and aggregate ScaleEvents
type ScaleAlertAggregator struct {
	logger        zerolog.Logger
	subscriptions []chan sokar.ScaleEvent

	// emitters is a list of components that are able to provide/ emit ScaleAlert's.
	emitters []ScaleAlertEmitter

	// stopChan is a channel used to signal teardown/ stop
	stopChan chan struct{}

	// weightMap contains a mapping of a ScalingAlert (specified by its name)
	// to a weight. A weight is defined in value per second.
	// This means the given weight is applied each second to the aggregated
	// scale counter of the ScaleAlertAggregator.
	// As soon as the scale counter exceeds the up-scale threshold
	// or underflows the down-scale threshold a scaling will be initiated.
	// Thus the higher (for up-scaling) the lower (for down-scaling) the weights
	// are the faster is the actual scale executed.
	weightMap ScaleAlertWeightMap

	// noAlertScaleDamping is a value that is applied each time no ScalingAlert is firing
	// (neither down nor upscaling). The value is used to move the
	// scaleCounter towards 0.
	// I.e. scaleCounter = scaleCounter + (sign(scaleCounter) -1) * noAlertScaleDamping
	noAlertScaleDamping float32

	scaleCounterGradient helper.LatestGradient

	// scaleCounter is used to decide wether to scale up/ down or wait.
	// The scaleCounter is a value representing the aggregated weighted scaling alerts.
	// To calculate this value the active scaling alerts (represented by their weight) are aggregated
	// by the ScaleAlertAggregator.
	scaleCounter scaleCounter

	// upScalingThreshold is the threshold that is used to trigger an upscaling event.
	// In case the scaleCounter is higher than this threshold, the upscaling event will be triggered.
	upScalingThreshold float32

	// downScalingThreshold is the threshold that is used to trigger an downscaling event.
	// In case the scaleCounter is lower than this threshold, the downscaling event will be triggered.
	downScalingThreshold float32

	// This map represents the ScaleAlerts which are currently known.
	// They where obtained through the different ScaleAlertEmitters
	scaleAlertPool ScaleAlertPool

	// aggregationCycle is the frequency the ScaleAlertAggregator will evaluate and aggregate the state
	// of the received ScaleAlert's
	aggregationCycle time.Duration

	// evaluationPeriodFactor is used to calculate the evaluation period.
	// The evaluation period is the period that is regarded for calculating the scaleCounterGradient/ scaleFactor.
	// Only the changes of the scaleCounter within this period/ window are regarded for scaleFactor calculation.
	// The period is a multiple of the aggregationCycle thus it is calculated by:
	// evaluationPeriod = aggregationCycle * evaluationPeriodFactor
	evaluationPeriodFactor uint

	// cleanupCycle is the frequency the ScaleAlertAggregator will cleanup/ remove
	// expired ScaleAlerts
	cleanupCycle time.Duration
}

// Config configuration for the ScaleAlertAggregator
type Config struct {
	Logger zerolog.Logger
}

// New creates a instance of the ScaleAlertAggregator
func (cfg Config) New(emitters []ScaleAlertEmitter) *ScaleAlertAggregator {
	return &ScaleAlertAggregator{
		logger:                 cfg.Logger,
		emitters:               emitters,
		stopChan:               make(chan struct{}, 1),
		scaleAlertPool:         NewScaleAlertPool(),
		aggregationCycle:       time.Millisecond * 2000,
		evaluationPeriodFactor: 10,
		cleanupCycle:           time.Second * 10,
		weightMap:              map[string]float32{"AlertA": 2.0, "AlertB": -1, "AlertC": -2},
		noAlertScaleDamping:    1.0,
		upScalingThreshold:     20.0,
		downScalingThreshold:   -20.0,
		scaleCounter:           newScaleCounter(),
		scaleCounterGradient:   helper.LatestGradient{Value: 0, Timestamp: time.Now()},
	}
}

// ScaleAlertWeightMap maps an alert name to its weight
type ScaleAlertWeightMap map[string]float32
