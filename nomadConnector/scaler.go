package nomadConnector

import (
	"fmt"
	"time"

	nomadApi "github.com/hashicorp/nomad/api"
	nomadstructs "github.com/hashicorp/nomad/nomad/structs"
)

const (
	deploymentTimeOut = 15 * time.Minute
	evaluationTimeOut = 30 * time.Second
)

// queryOptions sets sokars default QueryOptions for making GET calls to
// the nomad API.
func (nc *connectorImpl) queryOptions() (queryOptions *nomadApi.QueryOptions) {
	return &nomadApi.QueryOptions{AllowStale: true}
}

func (nc *connectorImpl) ScaleBy(amount int) error {
	nc.log.Info().Str("job", nc.jobName).Int("amount", amount).Msg("Scaling ...")

	// In order to scale the job, we need information on the current status of the
	// running job from Nomad.
	jobInfo, queryMeta, err := nc.nomad.Jobs().Info(nc.jobName, nc.queryOptions())
	nc.log.Debug().Uint64("LastIndex", queryMeta.LastIndex).Msg("QueryMeta: ")

	if err != nil {
		nc.log.Error().Err(err).Msg("Unable to determine job info")
		return err
	}

	// Use the current task count in order to determine whether or not a scaling
	// event will violate the min/max job policy.
	for i, _ := range jobInfo.TaskGroups {
		count := *jobInfo.TaskGroups[i].Count

		//if group.ScaleDirection == ScalingDirectionOut && *taskGroup.Count >= group.Max ||
		//	group.ScaleDirection == ScalingDirectionIn && *taskGroup.Count <= group.Min {
		//	logging.Debug("client/job_scaling: scale %v not permitted due to constraints on job \"%v\" and group \"%v\"",
		//		group.ScaleDirection, *jobInfo.ID, group.GroupName)
		//	return
		//}

		//logging.Info("client/job_scaling: scale %v will now be initiated against job \"%v\" and group \"%v\"",
		//	group.ScaleDirection, jobName, group.GroupName)

		// Depending on the scaling direction decrement/incrament the count;
		// currently replicator only supports addition/subtraction of 1.
		//if *taskGroup.Name == group.GroupName && group.ScaleDirection == ScalingDirectionOut {
		//	*jobResp.TaskGroups[i].Count++
		//	state.ScaleOutRequests++
		//}
		//
		//if *taskGroup.Name == group.GroupName && group.ScaleDirection == ScalingDirectionIn {
		//	*jobResp.TaskGroups[i].Count--
		//	state.ScaleInRequests++
		//}

		*jobInfo.TaskGroups[i].Count = count + amount
	}

	// Submit the job to the Register API endpoint with the altered count number
	// and check that no error is returned.
	jobRegisterResponse, writeMeta, err := nc.nomad.Jobs().Register(jobInfo, &nomadApi.WriteOptions{})
	nc.log.Debug().Uint64("LastIndex", writeMeta.LastIndex).Msg("WriteMeta: ")

	if err != nil {
		nc.log.Error().Err(err).Msg("Unable to scale")
		return err
	}

	//// Track the scaling submission time.
	//state.LastScalingEvent = time.Now()
	//if err != nil {
	//	logging.Error("client/job_scaling: issue submitting job %s for scaling action: %v", jobName, err)
	//	return
	//}
	//
	//// Setup our metric scaling direction namespace.
	//m := fmt.Sprintf("scale_%s", strings.ToLower(group.ScaleDirection))
	//
	err = nc.waitForScaleCompletion(jobRegisterResponse.EvalID, 15*time.Minute)
	if err != nil {
		nc.log.Error().Err(err).Msg("Failed scaling")
	}

	//if !success {
	//	metrics.IncrCounter([]string{"job", jobName, group.GroupName, m, "failure"}, 1)
	//	state.FailureCount++
	//
	//	return
	//}
	//
	//metrics.IncrCounter([]string{"job", jobName, group.GroupName, m, "success"}, 1)
	//logging.Info("client/job_scaling: scaling of job \"%v\" and group \"%v\" successfully completed",
	//	jobName, group.GroupName)

	nc.log.Info().Str("job", nc.jobName).Int("amount", amount).Msg("Scaling ...done")
	return nil
}

// waitForScaleCompletion checks if the deployment forced by the scale-event was successful or not.
func (nc *connectorImpl) waitForScaleCompletion(evalID string, timeout time.Duration) error {

	deplID, err := nc.getDeploymentID(evalID, 30*time.Second)
	if err != nil {
		return fmt.Errorf("Failed to retrieve deployment ID for evaluation %s.", evalID)
	}

	// Retry/ poll nomad each 500ms
	pollTicker := time.NewTicker(500 * time.Millisecond)
	defer pollTicker.Stop()

	deploymentTimeOut := time.After(timeout)

	queryOpt := &nomadApi.QueryOptions{WaitIndex: 1, AllowStale: true}

	for {
		select {

		// Timeout reached
		case <-deploymentTimeOut:
			return fmt.Errorf("Deployment (%s) timed out after %v", deplID, timeout)

		// Poll
		case <-pollTicker.C:
			deployment, queryMeta, err := nc.nomad.Deployments().Info(deplID, queryOpt)
			if err != nil {
				return err
			}

			// Wait/ redo until the waitIndex was transcended
			// It makes no sense to evaluate results earlier
			if queryMeta.LastIndex <= queryOpt.WaitIndex {
				continue
			}
			queryOpt.WaitIndex = queryMeta.LastIndex

			// Check the deployment status.
			if deployment.Status == nomadstructs.DeploymentStatusSuccessful {
				return nil
			} else if deployment.Status == nomadstructs.DeploymentStatusRunning {
				nc.log.Debug().Str("DeplID", deplID).Msg("Deployment still in progress.")
				continue
			} else {
				return fmt.Errorf("Deployment (%s) failed with status %s", deplID, deployment.Status)
			}
		}
	}
}

// getDeploymentID obtains the deployment ID of the given evaluation denoted by the evalID.
// Internally nomad is polled as long as the deployment ID was obtained successfully or
// the given timeout was reached.s
func (nc *connectorImpl) getDeploymentID(evalID string, timeout time.Duration) (depID string, err error) {

	// retry polling the nomad api until the deployment id was obtained successfully
	// or the evaluationTimeout was reached.
	pollTicker := time.NewTicker(time.Millisecond * 500)
	defer pollTicker.Stop()

	evaluationTimeout := time.After(timeout)

	for {
		select {

		// Timout Reached
		case <-evaluationTimeout:
			return depID, fmt.Errorf("EvaluationTimeout reached while trying to retrieve the "+
				"deployment ID for evaluation %v", evalID)

		// Retry
		case <-pollTicker.C:
			evaluation, _, err := nc.nomad.Evaluations().Info(evalID, nil)

			if err != nil {
				nc.log.Error().Str("EvalID", evalID).Err(err).Msg("Error while retrieving the deployment ID")
				continue
			}

			if evaluation.DeploymentID == "" {
				nc.log.Debug().Str("EvalID", evalID).Msg("Received deployment ID was empty. Will retry.")
				continue
			}

			nc.log.Debug().Str("EvalID", evalID).Str("DeplID", evaluation.DeploymentID).Msg("Received deployment ID.")

			return evaluation.DeploymentID, nil
		}
	}
}
