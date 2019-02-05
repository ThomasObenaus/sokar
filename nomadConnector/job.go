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

func (nc *connectorImpl) GetJobCount(jobname string) (uint, error) {
	jobInfo, err := nc.getJobInfo(jobname)

	if err != nil {
		return 0, err
	}

	// HACK: To unify the multiple groups with we take the job with max count.

	var count int
	for _, taskGroup := range jobInfo.TaskGroups {
		if *taskGroup.Count > count {
			count = *taskGroup.Count
		}
	}

	return uint(count), nil
}

func (nc *connectorImpl) SetJobCount(jobname string, count uint) error {
	nc.log.Info().Str("job", jobname).Msgf("Adjust job count of %s (including all groups) to %d.", jobname, count)

	// obtain current status about the job
	jobInfo, err := nc.getJobInfo(jobname)

	if err != nil {
		return err
	}

	// Use the current task count in order to determine whether or not a scaling
	// event will violate the min/max job policy.
	for _, taskGroup := range jobInfo.TaskGroups {
		nc.log.Info().Str("job", jobname).Str("grp", *taskGroup.Name).Msgf("Adjust count of group from %d to %d.", *taskGroup.Count, count)
		*taskGroup.Count = int(count)
	}

	// Submit the job to the Register API endpoint with the altered count number
	// and check that no error is returned.
	jobRegisterResponse, _, err := nc.jobsIF.Register(jobInfo, &nomadApi.WriteOptions{})

	if err != nil {
		nc.log.Error().Err(err).Str("Job", jobname).Msg("Unable to scale")
		return err
	}

	nc.log.Info().Str("job", jobname).Msg("Deployment issued, waiting for completion ... ")

	err = nc.waitForDeploymentConfirmation(jobRegisterResponse.EvalID, nc.deploymentTimeOut)

	if err != nil {
		nc.log.Error().Err(err).Msg("Deployment failed")
		return err
	}

	nc.log.Info().Str("job", jobname).Msg("Deployment issued, waiting for completion ... done")

	return nil
}
