package nomadConnector

import (
	"time"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
)

type connectorImpl struct {
	log zerolog.Logger

	// This is the object for interacting with nomad
	nomad *nomadApi.Client
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

func (nc *connectorImpl) SetJobCount(jobname string, count int) error {
	nc.log.Info().Str("job", jobname).Msgf("Adjust job count of %s (including all groups) to %d.", jobname, count)

	// In order to scale the job, we need information on the current status of the
	// running job from Nomad.
	jobInfo, _, err := nc.nomad.Jobs().Info(jobname, nc.defaultQueryOptions())

	if err != nil {
		nc.log.Error().Err(err).Msg("Unable to determine job info")
		return err
	}

	// Use the current task count in order to determine whether or not a scaling
	// event will violate the min/max job policy.
	for _, taskGroup := range jobInfo.TaskGroups {
		nc.log.Info().Str("job", jobname).Str("grp", *taskGroup.Name).Msgf("Adjust count of group from %d to %d.", *taskGroup.Count, count)
		*taskGroup.Count = count
	}

	// Submit the job to the Register API endpoint with the altered count number
	// and check that no error is returned.
	jobRegisterResponse, _, err := nc.nomad.Jobs().Register(jobInfo, &nomadApi.WriteOptions{})

	if err != nil {
		nc.log.Error().Err(err).Str("Job", jobname).Msg("Unable to scale")
		return err
	}

	nc.log.Info().Str("job", jobname).Msg("Deployment issued, waiting for completion ... ")

	err = nc.waitForDeploymentConfirmation(jobRegisterResponse.EvalID, 15*time.Minute)

	if err != nil {
		nc.log.Error().Err(err).Msg("Deployment failed")
	}

	nc.log.Info().Str("job", jobname).Msg("Deployment issued, waiting for completion ... done")

	return nil
}
