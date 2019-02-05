package nomadConnector

import (
	"fmt"
	"time"

	nomadApi "github.com/hashicorp/nomad/api"
	nomadstructs "github.com/hashicorp/nomad/nomad/structs"
)

func (nc *connectorImpl) printDeploymentProgress(deplID string, deployment *nomadApi.Deployment) {
	nc.log.Debug().Str("DeplID", deplID).Msgf("Deployment still in progress (%s).", deployment.StatusDescription)
	for tgName, deplState := range deployment.TaskGroups {
		nc.log.Debug().Str("DeplID", deplID).Msgf("taskGroup=%s, Allocs: desired=%d,placed=%d,healthy=%d,unhealthy=%d", tgName, deplState.DesiredTotal, deplState.PlacedAllocs, deplState.HealthyAllocs, deplState.UnhealthyAllocs)
	}
}

// waitForDeploymentConfirmation checks if the deployment forced by the scale-event was successful or not.
func (nc *connectorImpl) waitForDeploymentConfirmation(evalID string, timeout time.Duration) error {

	deplID, err := nc.getDeploymentID(evalID, nc.evaluationTimeOut)
	if err != nil {
		return fmt.Errorf("Failed to retrieve deployment ID for evaluation %s: %s", evalID, err.Error())
	}

	// Retry/ poll nomad each 500ms
	pollTicker := time.NewTicker(500 * time.Millisecond)
	defer pollTicker.Stop()

	deploymentTimeOut := time.After(timeout)

	queryOpt := &nomadApi.QueryOptions{WaitIndex: 1, AllowStale: true, WaitTime: time.Second * 15}

	for {
		select {

		// Timeout reached
		case <-deploymentTimeOut:
			return fmt.Errorf("Deployment (%s) timed out after %v", deplID, timeout)

		// Poll
		case <-pollTicker.C:
			deployment, queryMeta, err := nc.deploymentIF.Info(deplID, queryOpt)
			if err != nil {
				return err
			}

			if deployment == nil || queryMeta == nil {
				return fmt.Errorf("Got nil while querying for deployment %s", deplID)
			}

			// Wait/ redo until the waitIndex was transcended
			// It makes no sense to evaluate results earlier
			if queryMeta.LastIndex <= queryOpt.WaitIndex {
				nc.log.Warn().Str("DeplID", deplID).Msgf("Waitindex not exceeded yet (lastIdx=%d, waitIdx=%d). Probably resources are exhausted.", queryMeta.LastIndex, queryOpt.WaitIndex)
				nc.printDeploymentProgress(deplID, deployment)
				continue
			}

			queryOpt.WaitIndex = queryMeta.LastIndex

			// Check the deployment status.
			if deployment.Status == nomadstructs.DeploymentStatusSuccessful {
				return nil
			} else if deployment.Status == nomadstructs.DeploymentStatusRunning {
				nc.printDeploymentProgress(deplID, deployment)
				continue
			} else {
				return fmt.Errorf("Deployment (%s) failed with status %s (%s).", deplID, deployment.Status, deployment.StatusDescription)
			}
		}
	}
}

// getDeploymentID obtains the deployment ID of the given evaluation denoted by the evalID.
// Internally nomad is polled as long as the deployment ID was obtained successfully or
// the given timeout was reached.s
func (nc *connectorImpl) getDeploymentID(evalID string, timeout time.Duration) (depID string, err error) {

	evalIf := nc.evalIF
	if evalIf == nil {
		return "", fmt.Errorf("Nomad Evaluations() interface is missing")
	}

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
			evaluation, _, err := nc.evalIF.Info(evalID, nil)

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
