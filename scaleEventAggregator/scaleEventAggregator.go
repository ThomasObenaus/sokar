package scaleEventAggregator

import (
	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/sokar"
)

type ScaleEventAggregator struct {
	logger        zerolog.Logger
	subscriptions []chan sokar.ScaleEvent
}

type Config struct {
	Logger zerolog.Logger
}

func (cfg Config) New() *ScaleEventAggregator {
	return &ScaleEventAggregator{
		logger: cfg.Logger,
	}
}
