package nomadWorker

import (
	"github.com/rs/zerolog"
)

// Connector is a object that allows to interact with nomad worker
type Connector struct {
	log zerolog.Logger
}

// Config contains the main configuration for the nomad worker connector
type Config struct {
	Logger zerolog.Logger
}

// New creates a new nomad connector
func (cfg *Config) New() (*Connector, error) {

	nc := &Connector{
		log: cfg.Logger,
	}

	cfg.Logger.Info().Msg("Setting up nomad worker connector ... done")
	return nc, nil
}
