package nomadConnector

import (
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
)

type connectorImpl struct {
	jobName string
	log     zerolog.Logger

	// This is the object for interacting with nomad
	nomad *nomadApi.Client
}
