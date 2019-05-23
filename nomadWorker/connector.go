package nomadWorker

import (
	"github.com/rs/zerolog"
)

// Connector is a object that allows to interact with nomad worker
type Connector struct {
	log zerolog.Logger

	currentCount uint
}

// Config contains the main configuration for the nomad worker connector
type Config struct {
	Logger zerolog.Logger
}

// New creates a new nomad connector
func (cfg *Config) New() (*Connector, error) {

	nc := &Connector{
		log: cfg.Logger,
		// HACK: Set it to 100 for now to ensure at startup that a scale is possible (i.e. with a value 0 a initial downscale would be ignored)
		currentCount: 100,
	}

	cfg.Logger.Info().Msg("Setting up nomad worker connector ... done")
	return nc, nil
}
