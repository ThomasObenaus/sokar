package nomadConnector

import (
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
)

type Config struct {
	JobName            string
	NomadServerAddress string
}

// Connector defines the interface of the component being able to communicate with nomad
type Connector interface {
	ScaleBy(amount int) error
}

// New creates a new nomad connector
func (cfg *Config) New(logger zerolog.Logger) (Connector, error) {
	logger.Info().Str("srvAddr", cfg.NomadServerAddress).Msg("Setting up nomad connector ...")

	// config needed to set up a nomad api client
	config := nomadApi.DefaultConfig()
	config.Address = cfg.NomadServerAddress
	//config.SecretID = token
	//config.TLSConfig.TLSServerName = tls_server_name

	client, err := nomadApi.NewClient(config)
	if err != nil {
		return nil, err
	}

	// issue test query to find out if the connection to nomad works
	peers, err := client.Status().Peers()
	if err != nil {
		return nil, err
	}

	nc := &connectorImpl{
		jobName: cfg.JobName,
		log:     logger,
		nomad:   client,
	}

	logger.Info().Str("srvAddr", cfg.NomadServerAddress).Int("#peers", len(peers)).Msg("Setting up nomad connector ... done")
	return nc, nil
}
