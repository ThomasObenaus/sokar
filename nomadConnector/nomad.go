package nomadConnector

import (
	"time"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
)

type connectorImpl struct {
	log zerolog.Logger

	jobsIF       NomadJobs
	deploymentIF NomadDeployments
	evalIF       NomadEvaluations
}

const (
	deploymentTimeOut = 15 * time.Minute
	evaluationTimeOut = 30 * time.Second
)

// defaultQueryOptions sets sokars default QueryOptions for making GET calls to
// the nomad API.
func (nc *connectorImpl) defaultQueryOptions() (queryOptions *nomadApi.QueryOptions) {
	return &nomadApi.QueryOptions{AllowStale: true}
}
