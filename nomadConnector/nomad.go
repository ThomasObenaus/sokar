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

	deploymentTimeOut time.Duration
	evaluationTimeOut time.Duration
}

// defaultQueryOptions sets sokars default QueryOptions for making GET calls to
// the nomad API.
func (nc *connectorImpl) defaultQueryOptions() (queryOptions *nomadApi.QueryOptions) {
	return &nomadApi.QueryOptions{AllowStale: true}
}
