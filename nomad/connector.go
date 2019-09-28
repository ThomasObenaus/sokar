package nomad

import (
	"fmt"
	"time"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
)

// Connector is a object that allows to interact with nomad
type Connector struct {
	log zerolog.Logger

	// Interfaces needed to interact with nomad
	jobsIF       Jobs
	deploymentIF Deployments
	evalIF       Evaluations

	deploymentTimeOut time.Duration
	evaluationTimeOut time.Duration
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

// Option represents an option for the Connector
type Option func(conn *Connector)

// WithLogger adds a configured Logger to the Connector
func WithLogger(logger zerolog.Logger) Option {
	return func(conn *Connector) {
		conn.log = logger
	}
}

// WithDeploymentTimeOut sets the timeout that should be regarded during deployments
func WithDeploymentTimeOut(timeout time.Duration) Option {
	return func(conn *Connector) {
		conn.deploymentTimeOut = timeout
	}
}

// WithEvaluationTimeOut sets the timeout that should be regarded during evaluations
func WithEvaluationTimeOut(timeout time.Duration) Option {
	return func(conn *Connector) {
		conn.evaluationTimeOut = timeout
	}
}

// New creates a new nomad connector
func New(nomadServerAddress string, options ...Option) (*Connector, error) {

	if len(nomadServerAddress) == 0 {
		return nil, fmt.Errorf("required configuration 'nomadServerAddress' is missing/ empty")
	}

	nomadConnector := Connector{
		deploymentTimeOut: 1 * time.Minute,
		evaluationTimeOut: 30 * time.Second,
	}
	// apply the options
	for _, opt := range options {
		opt(&nomadConnector)
	}

	nomadConnector.log.Info().Str("srvAddr", nomadServerAddress).Msg("Setting up nomad connector ...")

	// config needed to set up a nomad api client
	config := nomadApi.DefaultConfig()
	config.Address = nomadServerAddress
	//config.SecretID = token
	//config.TLSConfig.TLSServerName = tls_server_name

	client, err := nomadApi.NewClient(config)
	if err != nil {
		return nil, err
	}

	nomadConnector.jobsIF = client.Jobs()
	nomadConnector.deploymentIF = client.Deployments()
	nomadConnector.evalIF = client.Evaluations()

	nomadConnector.log.Info().Str("srvAddr", nomadServerAddress).Msg("Setting up nomad connector ... done")
	return &nomadConnector, nil
}

func (c *Connector) String() string {
	return "Nomad-Job"
}
