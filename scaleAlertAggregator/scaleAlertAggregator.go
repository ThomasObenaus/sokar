package scaleAlertAggregator

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/helper"
	sokar "github.com/thomasobenaus/sokar/sokar/iface"
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
	scaleCounter float32

	// upScalingThreshold is the threshold that is used to trigger an upscaling event.
	// In case the scaleCounter is higher than this threshold, the upscaling event will be triggered.
	upScalingThreshold float32

	// downScalingThreshold is the threshold that is used to trigger an downscaling event.
	// In case the scaleCounter is lower than this threshold, the downscaling event will be triggered.
	downScalingThreshold float32

	// This map represents the ScaleAlerts which are currently known.
	// They where obtained through the different ScaleAlertEmitters
	scaleAlertPool ScaleAlertPool

	// evaluationCycle is the frequency the ScaleAlertAggregator will evaluate and aggregate the state
	// of the received ScaleAlert's
	evaluationCycle time.Duration

	// evaluationPeriodFactor is used to calculate the evaluation period.
	// The evaluation period is the period that is regarded for calculating the scaleCounterGradient/ scaleFactor.
	// Only the changes of the scaleCounter within this period/ window are regarded for scaleFactor calculation.
	// The period is a multiple of the evaluationCycle thus it is calculated by:
	// evaluationPeriod = evaluationCycle * evaluationPeriodFactor
	evaluationPeriodFactor uint

	evaluationCounter uint

	// cleanupCycle is the frequency the ScaleAlertAggregator will cleanup/ remove
	// expired ScaleAlerts
	cleanupCycle time.Duration

	// metrics is a collection of metrics used by the ScaleAlertAggregator
	metrics Metrics
}

// Config configuration for the ScaleAlertAggregator
type Config struct {
	Logger zerolog.Logger

	WeightMap              ScaleAlertWeightMap
	NoAlertScaleDamping    float32
	UpScalingThreshold     float32
	DownScalingThreshold   float32
	EvaluationCycle        time.Duration
	EvaluationPeriodFactor uint
	CleanupCycle           time.Duration

	// AlertExpirationTime defines after which time an alert will be pruned if he did not
	// get updated again by the ScaleAlertEmitter, assuming that the alert is not relevant any more.
	AlertExpirationTime time.Duration
}

// NewDefaultConfig creates an empty default configuration
func NewDefaultConfig() Config {
	return Config{
		WeightMap:              make(ScaleAlertWeightMap),
		NoAlertScaleDamping:    1,
		UpScalingThreshold:     10,
		DownScalingThreshold:   -10,
		EvaluationCycle:        time.Second * 1,
		EvaluationPeriodFactor: 10,
		CleanupCycle:           time.Second * 60,
		AlertExpirationTime:    time.Minute * 10,
	}
}

// New creates a instance of the ScaleAlertAggregator
func (cfg Config) New(emitters []ScaleAlertEmitter, metrics Metrics) *ScaleAlertAggregator {
	return &ScaleAlertAggregator{
		logger:                 cfg.Logger,
		emitters:               emitters,
		stopChan:               make(chan struct{}, 1),
		scaleAlertPool:         NewScaleAlertPool(cfg.AlertExpirationTime),
		evaluationCycle:        cfg.EvaluationCycle,
		evaluationPeriodFactor: cfg.EvaluationPeriodFactor,
		cleanupCycle:           cfg.CleanupCycle,
		weightMap:              cfg.WeightMap,
		noAlertScaleDamping:    cfg.NoAlertScaleDamping,
		upScalingThreshold:     cfg.UpScalingThreshold,
		downScalingThreshold:   cfg.DownScalingThreshold,
		scaleCounter:           0,
		scaleCounterGradient:   helper.LatestGradient{Value: 0, Timestamp: time.Now()},
		evaluationCounter:      0,
		metrics:                metrics,
	}
}

// ScaleAlertWeightMap maps an alert name to its weight
type ScaleAlertWeightMap map[string]float32
