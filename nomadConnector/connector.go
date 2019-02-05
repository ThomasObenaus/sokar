package nomadConnector

import (
	"fmt"
	"os"
	"time"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Connector defines the interface of the component being able to communicate with nomad
type Connector interface {
	SetJobCount(jobname string, count uint) error
	GetJobCount(jobname string) (uint, error)
}

// Config contains the main configuration for the nomad connector
type Config struct {
	NomadServerAddress string
	Logger             zerolog.Logger

	// DeploymentTimeOut reflects the timeout sokar will wait (at max) for a deployment to be applied.
	DeploymentTimeOut time.Duration
	// EvaluationTimeOut reflects the timeout sokar will wait (at max) for gathering information about evaluations.
	EvaluationTimeOut time.Duration
}

// NewDefaultConfig returns a good default configuration for the nomad connector
func NewDefaultConfig(nomadServerAddress string) Config {
	return Config{
		NomadServerAddress: nomadServerAddress,
		Logger:             log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Str("logger", "sokar").Logger(),
		DeploymentTimeOut:  1 * time.Minute,
		EvaluationTimeOut:  30 * time.Second,
	}
}

// New creates a new nomad connector
func (cfg *Config) New() (Connector, error) {

	if len(cfg.NomadServerAddress) == 0 {
		return nil, fmt.Errorf("Required configuration 'NomadServerAddress' is missing.")
	}

	cfg.Logger.Info().Str("srvAddr", cfg.NomadServerAddress).Msg("Setting up nomad connector ...")

	// config needed to set up a nomad api client
	config := nomadApi.DefaultConfig()
	config.Address = cfg.NomadServerAddress
	//config.SecretID = token
	//config.TLSConfig.TLSServerName = tls_server_name

	client, err := nomadApi.NewClient(config)
	if err != nil {
		return nil, err
	}

	nc := &connectorImpl{
		log:               cfg.Logger,
		jobsIF:            client.Jobs(),
		deploymentIF:      client.Deployments(),
		evalIF:            client.Evaluations(),
		deploymentTimeOut: cfg.DeploymentTimeOut,
		evaluationTimeOut: cfg.EvaluationTimeOut,
	}

	cfg.Logger.Info().Str("srvAddr", cfg.NomadServerAddress).Msg("Setting up nomad connector ... done")
	return nc, nil
}
