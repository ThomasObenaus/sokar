package nomadConnector

import (
	"github.com/rs/zerolog"
)

type Config struct {
	JobName string
}

// Connector defines the interface of the component being able to communicate with nomad
type Connector interface {
	ScaleBy(amount int) error
}

// New creates a new nomad connector
func (cfg *Config) New(logger zerolog.Logger) Connector {
	nc := &connectorImpl{
		jobName: cfg.JobName,
		log:     logger,
	}

	return nc
}
