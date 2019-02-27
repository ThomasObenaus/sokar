package scaleEventAggregator

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/sokar"
)

// ScaleEventAggregator is a component that is responsible to gather and aggregate ScaleEvents
type ScaleEventAggregator struct {
	logger        zerolog.Logger
	subscriptions []chan sokar.ScaleEvent

	// A list of components that are able to provide ScalAlerts.
	receivers []ScaleAlertReceiver

	// channel used to signal teardown/ stop
	stopChan chan struct{}

	// The weightMap contains a mapping of a ScalingAlert (specified by its name)
	// to a weight. A weight is defined in value per second.
	// This means the given weight is applied each second to the aggregated
	// scale counter of the ScaleEventAggregator.
	// As soon as the scale counter exceeds the up-scale threshold
	// or underflows the down-scale threshold a scaling will be initiated.
	// Thus the higher (for up-scaling) the lower (for down-scaling) the weights
	// are the faster is the actual scale executed.
	weightMap ScaleAlertWeightMap

	// This value is applied each time no ScalingAlert is firing
	// (neither down nor upscaling). The value is used to move the
	// scaleCounter towards 0.
	// I.e. scaleCounter = scaleCounter + (sign(scaleCounter) -1) * noAlertScaleDamping
	noAlertScaleDamping float32

	// scaleCounter is used to decide wether to scale up/ down or wait.
	// The scaleCounter is a value representing the aggregated weighted scaling alerts.
	// To calculate this value the active scaling alerts (represented by their weight) are aggregated
	// by the ScaleEventAggregator.
	scaleCounter float32

	// upScalingThreshold is the threshold that is used to trigger an upscaling event.
	// In case the scaleCounter is higher than this threshold, the upscaling event will be triggered.
	upScalingThreshold float32

	// downScalingThreshold is the threshold that is used to trigger an downscaling event.
	// In case the scaleCounter is lower than this threshold, the downscaling event will be triggered.
	downScalingThreshold float32

	// This map represents the ScaleAlerts which are currently known.
	// They where obtained through the different ScaleAlertReceivers
	scaleAlertPool ScaleAlertPool

	// The frequency the ScaleEventAggregator will evaluate and aggregate the state
	// of the received ScaleAlert's
	aggregationCycle time.Duration

	// The frequency the ScaleEventAggregator will cleanup/ remove
	// expired ScaleAlerts
	cleanupCycle time.Duration
}

// Config configuration for the ScaleEventAggregator
type Config struct {
	Logger zerolog.Logger
}

// New creates a instance of the ScaleEventAggregator
func (cfg Config) New(receivers []ScaleAlertReceiver) *ScaleEventAggregator {
	return &ScaleEventAggregator{
		logger:               cfg.Logger,
		receivers:            receivers,
		stopChan:             make(chan struct{}, 1),
		scaleAlertPool:       NewScaleAlertPool(),
		aggregationCycle:     time.Millisecond * 2000,
		cleanupCycle:         time.Second * 10,
		weightMap:            map[string]float32{"AlertA": 2.0, "AlertB": -1},
		noAlertScaleDamping:  1.0,
		upScalingThreshold:   5.0,
		downScalingThreshold: -5.0,
		scaleCounter:         0,
	}
}

// ScaleAlertWeightMap maps an alert name to its weight
type ScaleAlertWeightMap map[string]float32
