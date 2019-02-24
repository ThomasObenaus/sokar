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
	receivers     []ScaleAlertReceiver

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

	scaleCounter         float32
	upScalingThreshold   float32
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
		upScalingThreshold:   5.0,
		downScalingThreshold: -5.0,
		scaleCounter:         0,
	}
}

// ScaleAlertWeightMap maps an alert name to its weight
type ScaleAlertWeightMap map[string]float32
