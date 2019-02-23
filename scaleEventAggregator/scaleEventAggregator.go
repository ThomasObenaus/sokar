package scaleEventAggregator

import (
	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/sokar"
)

// ScaleEventAggregator is a component that is responsible to gather and aggregate ScaleEvents
type ScaleEventAggregator struct {
	logger        zerolog.Logger
	subscriptions []chan sokar.ScaleEvent
}

// Config configuration for the ScaleEventAggregator
type Config struct {
	Logger zerolog.Logger
}

// New creates a instance of the ScaleEventAggregator
func (cfg Config) New() *ScaleEventAggregator {
	return &ScaleEventAggregator{
		logger: cfg.Logger,
	}
}
