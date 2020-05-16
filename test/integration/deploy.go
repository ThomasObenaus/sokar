package main

import (
	"fmt"
	"time"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/pkg/errors"
	"github.com/thomasobenaus/sokar/nomad/structs"
)

type deployerImpl struct {

	// Interfaces needed to interact with nomad
	jobsIF       *nomadApi.Jobs
	deploymentIF *nomadApi.Deployments
	evalIF       *nomadApi.Evaluations

	deploymentTimeout time.Duration
	evaluationTimeOut time.Duration
}

func NewDeployer(nomadServerAddress string) (*deployerImpl, error) {

	config := nomadApi.DefaultConfig()
	config.Address = nomadServerAddress
	client, err := nomadApi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &deployerImpl{
		jobsIF:            client.Jobs(),
		deploymentIF:      client.Deployments(),
		evalIF:            client.Evaluations(),
		deploymentTimeout: time.Second * 30,
		evaluationTimeOut: time.Second * 15,
	}, nil
}

func strToPtr(v string) *string {
	return &v
}

func intToPtr(v int) *int {
	return &v
}

func durToPtr(v time.Duration) *time.Duration {
	return &v
}

func (d *deployerImpl) Deploy() error {

	nwResource := nomadApi.NetworkResource{
		MBits:        intToPtr(10),
		DynamicPorts: []nomadApi.Port{nomadApi.Port{Label: "http"}},
	}
	resources := nomadApi.Resources{
		CPU:      intToPtr(100),
		MemoryMB: intToPtr(128),
		Networks: []*nomadApi.NetworkResource{&nwResource},
	}

	service1 := nomadApi.Service{
		Name:      "fail-service",
		PortLabel: "http",
		Checks: []nomadApi.ServiceCheck{nomadApi.ServiceCheck{
			PortLabel: "http",
			Type:      "http",
			Path:      "/health",
			Method:    "GET",
			Interval:  time.Second * 10,
			Timeout:   time.Second * 2,
		}},
	}

	task1 := nomadApi.Task{
		Name:   "fail-service",
		Driver: "docker",
		Config: map[string]interface{}{
			"image":    "thobe/fail_service:v0.1.0",
			"port_map": []map[string]int{map[string]int{"http": 8080}},
		},
		Resources: &resources,
		Services:  []*nomadApi.Service{&service1},
		Env:       map[string]string{"HEALTHY_FOR": "-1"},
	}

	tasks1 := []*nomadApi.Task{&task1}
	taskGroup1 := nomadApi.TaskGroup{
		Name:  strToPtr("fail-service-grp-A"),
		Tasks: tasks1,
		Count: intToPtr(2),
	}
	taskGroups1 := []*nomadApi.TaskGroup{&taskGroup1}

	updateStrategy := nomadApi.UpdateStrategy{
		Stagger:     durToPtr(time.Second * 5),
		MaxParallel: intToPtr(1),
	}
	jobInfo := &nomadApi.Job{
		ID:          strToPtr("fail-service"),
		Datacenters: []string{"testing"},
		TaskGroups:  taskGroups1,
		Type:        strToPtr("service"),
		Update:      &updateStrategy,
	}

	fmt.Printf("[deploy] Register job\n")
	jobRegisterResponse, _, err := d.jobsIF.Register(jobInfo, &nomadApi.WriteOptions{})
	fmt.Printf("[deploy] Job registered resp='%v'\n", jobRegisterResponse)

	if err != nil {
		return errors.Wrap(err, "Failed to register job for deployment")
	}

	err = d.waitForDeploymentConfirmation(jobRegisterResponse.EvalID, d.deploymentTimeout)

	if err != nil {
		return errors.Wrap(err, "Deployment failed")
	}

	//nc.log.Info().Str("job", jobname).Msg("Deployment issued, waiting for completion ... done")
	return nil
}

// waitForDeploymentConfirmation checks if the deployment forced by the scale-event was successful or not.
func (d *deployerImpl) waitForDeploymentConfirmation(evalID string, timeout time.Duration) error {
	fmt.Printf("[deploy] Get deployment id for evailid=%s\n", evalID)

	deplID, err := d.getDeploymentID(evalID, d.evaluationTimeOut)
	if err != nil {
		return fmt.Errorf("Failed to retrieve deployment ID for evaluation %s: %s", evalID, err.Error())
	}
	fmt.Printf("[deploy] Got deployment id=%s\n", deplID)

	// Retry/ poll nomad each 500ms
	pollTicker := time.NewTicker(500 * time.Millisecond)
	defer pollTicker.Stop()

	deploymentTimeOut := time.After(timeout)

	queryOpt := &nomadApi.QueryOptions{WaitIndex: 1, AllowStale: true, WaitTime: time.Second * 15}

	deploymentIF := d.deploymentIF
	for {
		select {

		// Timeout reached
		case <-deploymentTimeOut:
			return fmt.Errorf("Deployment (%s) timed out after %v", deplID, timeout)

		// Poll
		case <-pollTicker.C:
			deployment, queryMeta, err := deploymentIF.Info(deplID, queryOpt)
			if err != nil {
				return err
			}
			fmt.Printf("[deploy] Pending deployment: %v\n", *deployment)

			if deployment == nil || queryMeta == nil {
				return fmt.Errorf("Got nil while querying for deployment %s", deplID)
			}

			// Wait/ redo until the waitIndex was transcended
			// It makes no sense to evaluate results earlier
			if queryMeta.LastIndex <= queryOpt.WaitIndex {
				//nc.log.Warn().Str("DeplID", deplID).Msgf("Waitindex not exceeded yet (lastIdx=%d, waitIdx=%d). Probably resources are exhausted.", queryMeta.LastIndex, queryOpt.WaitIndex)
				//nc.printDeploymentProgress(deplID, deployment)
				continue
			}

			queryOpt.WaitIndex = queryMeta.LastIndex

			// Check the deployment status.
			if deployment.Status == structs.DeploymentStatusSuccessful {
				return nil
			} else if deployment.Status == structs.DeploymentStatusRunning {
				//nc.printDeploymentProgress(deplID, deployment)
				continue
			} else {
				return fmt.Errorf("Deployment (%s) failed with status %s (%s)", deplID, deployment.Status, deployment.StatusDescription)
			}
		}
	}
}

// getDeploymentID obtains the deployment ID of the given evaluation denoted by the evalID.
// Internally nomad is polled as long as the deployment ID was obtained successfully or
// the given timeout was reached.s
func (d *deployerImpl) getDeploymentID(evalID string, timeout time.Duration) (depID string, err error) {

	evalIf := d.evalIF
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
			evaluation, _, err := evalIf.Info(evalID, nil)

			if err != nil {
				//nc.log.Error().Str("EvalID", evalID).Err(err).Msg("Error while retrieving the deployment ID")
				continue
			}

			if evaluation.DeploymentID == "" {
				//nc.log.Debug().Str("EvalID", evalID).Msg("Received deployment ID was empty. Will retry.")
				continue
			}

			//nc.log.Debug().Str("EvalID", evalID).Str("DeplID", evaluation.DeploymentID).Msg("Received deployment ID.")

			return evaluation.DeploymentID, nil
		}
	}
}
