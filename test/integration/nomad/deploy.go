package nomad

import (
	"fmt"
	"testing"
	"time"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/pkg/errors"
	"github.com/thomasobenaus/sokar/nomad/structs"
	"github.com/thomasobenaus/sokar/test/integration/helper"
)

type deployerImpl struct {

	// Interfaces needed to interact with nomad
	jobsIF       *nomadApi.Jobs
	deploymentIF *nomadApi.Deployments
	evalIF       *nomadApi.Evaluations

	deploymentTimeout time.Duration
	evaluationTimeOut time.Duration

	// tstCtx is the testing context
	tstCtx *testing.T
}

func NewDeployer(t *testing.T, nomadServerAddress string) (*deployerImpl, error) {

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
		tstCtx:            t,
	}, nil
}

// NewJobDescription creates a new default job description that can be used to deploy a nomad job.
// Per default the type of that job is 'service'.
func NewJobDescription(jobName, datacenter, dockerImage string, count int, envVars map[string]string) *nomadApi.Job {
	nwResource := nomadApi.NetworkResource{
		MBits:        helper.IntToPtr(10),
		DynamicPorts: []nomadApi.Port{nomadApi.Port{Label: "http"}},
	}
	resources := nomadApi.Resources{
		CPU:      helper.IntToPtr(100),
		MemoryMB: helper.IntToPtr(128),
		Networks: []*nomadApi.NetworkResource{&nwResource},
	}

	service := nomadApi.Service{
		Name:      fmt.Sprintf("%s-service", jobName),
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

	task := nomadApi.Task{
		Name:   fmt.Sprintf("%s-task", jobName),
		Driver: "docker",
		Config: map[string]interface{}{
			"image":    dockerImage,
			"port_map": []map[string]int{map[string]int{"http": 8080}},
		},
		Resources: &resources,
		Services:  []*nomadApi.Service{&service},
		Env:       envVars,
	}

	tasks := []*nomadApi.Task{&task}
	taskGroup := nomadApi.TaskGroup{
		Name:  helper.StrToPtr(fmt.Sprintf("%s-grp", jobName)),
		Tasks: tasks,
		Count: helper.IntToPtr(count),
	}
	taskGroups := []*nomadApi.TaskGroup{&taskGroup}

	updateStrategy := nomadApi.UpdateStrategy{
		Stagger:     helper.DurToPtr(time.Second * 5),
		MaxParallel: helper.IntToPtr(1),
	}
	jobInfo := &nomadApi.Job{
		ID:          helper.StrToPtr(jobName),
		Datacenters: []string{datacenter},
		TaskGroups:  taskGroups,
		Type:        helper.StrToPtr("service"),
		Update:      &updateStrategy,
	}
	return jobInfo
}

func (d *deployerImpl) Deploy(job *nomadApi.Job) error {

	fmt.Printf("[deploy] Register job\n")
	jobRegisterResponse, _, err := d.jobsIF.Register(job, &nomadApi.WriteOptions{})
	fmt.Printf("[deploy] Job registered resp='%v'\n", jobRegisterResponse)

	if err != nil {
		return errors.Wrap(err, "Failed to register job for deployment")
	}

	err = d.waitForDeploymentConfirmation(jobRegisterResponse.EvalID, d.deploymentTimeout)

	if err != nil {
		return errors.Wrap(err, "Deployment failed")
	}

	fmt.Printf("[deploy] Deployment done\n")
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

			if deployment == nil || queryMeta == nil {
				return fmt.Errorf("Got nil while querying for deployment %s", deplID)
			}

			d.printDeploymentProgress(deplID, deployment)

			// Wait/ redo until the waitIndex was transcended
			// It makes no sense to evaluate results earlier
			if queryMeta.LastIndex <= queryOpt.WaitIndex {
				fmt.Printf("[deploy][WARN] Waitindex not exceeded yet (lastIdx=%d, waitIdx=%d). Probably resources are exhausted.", queryMeta.LastIndex, queryOpt.WaitIndex)
				d.printDeploymentProgress(deplID, deployment)
				continue
			}

			queryOpt.WaitIndex = queryMeta.LastIndex

			// Check the deployment status.
			if deployment.Status == structs.DeploymentStatusSuccessful {
				return nil
			} else if deployment.Status == structs.DeploymentStatusRunning {
				d.printDeploymentProgress(deplID, deployment)
				continue
			} else {
				return fmt.Errorf("Deployment (%s) failed with status %s (%s)", deplID, deployment.Status, deployment.StatusDescription)
			}
		}
	}
}

func (d *deployerImpl) printDeploymentProgress(deplID string, deployment *nomadApi.Deployment) {
	fmt.Printf("[deploy] Deployment progress info (%s)\n", deployment.StatusDescription)
	for tgName, deplState := range deployment.TaskGroups {
		perc := (float32(deplState.HealthyAllocs) / float32(deplState.DesiredTotal)) * 100.0
		fmt.Printf("[deploy] taskGroup=%s, depl=%.2f%%, Allocs: desired=%d,placed=%d,healthy=%d,unhealthy=%d\n", tgName, perc, deplState.DesiredTotal, deplState.PlacedAllocs, deplState.HealthyAllocs, deplState.UnhealthyAllocs)
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
				fmt.Printf("[deploy][ERR] Error while retrieving the deployment ID: %s", err.Error())
				continue
			}

			if evaluation.DeploymentID == "" {
				continue
			}
			return evaluation.DeploymentID, nil
		}
	}
}
