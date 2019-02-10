package scaleEventAggregator

import (
	"github.com/rs/zerolog"
)

type ScaleEventAggregator struct {
	logger zerolog.Logger
}

type Config struct {
	Logger zerolog.Logger
}

func (cfg Config) New() *ScaleEventAggregator {
	return &ScaleEventAggregator{
		logger: cfg.Logger,
	}
}
