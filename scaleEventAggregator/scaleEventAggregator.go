package scaleEventAggregator

import (
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

	scaleFactorMap AlertToScaleFactorMap

	// This map represents the ScaleAlerts which are currently known.
	// They where obtained through the different ScaleAlertReceivers
	scaleAlertPool ScaleAlertPool
}

// Config configuration for the ScaleEventAggregator
type Config struct {
	Logger zerolog.Logger
}

// New creates a instance of the ScaleEventAggregator
func (cfg Config) New(receivers []ScaleAlertReceiver) *ScaleEventAggregator {
	return &ScaleEventAggregator{
		logger:         cfg.Logger,
		receivers:      receivers,
		stopChan:       make(chan struct{}, 1),
		scaleFactorMap: map[string]float32{"AlertA": -1.0, "AlertB": 2},
		scaleAlertPool: NewScaleAlertPool(),
	}
}

type AlertToScaleFactorMap map[string]float32
