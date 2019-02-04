package nomadConnector

import (
	"fmt"

	nomadApi "github.com/hashicorp/nomad/api"
)

func (nc *connectorImpl) getJobInfo(jobname string) (*nomadApi.Job, error) {
	jobs := nc.jobsIF

	if jobs == nil {
		return nil, fmt.Errorf("Nomad Jobs() interface is missing")
	}

	// In order to scale the job, we need information on the current status of the
	// running job from Nomad.
	jobInfo, _, err := jobs.Info(jobname, nc.defaultQueryOptions())

	if err != nil {
		nc.log.Error().Err(err).Msg("Unable to determine job info")
		return nil, err
	}

	return jobInfo, nil
}
