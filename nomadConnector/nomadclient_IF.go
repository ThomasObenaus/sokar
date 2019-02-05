package nomadConnector

import (
	nomadApi "github.com/hashicorp/nomad/api"
)

// NomadJobs represents the interface for interacting
// with nomads over the jobs end-point
type NomadJobs interface {
	Info(jobID string, q *nomadApi.QueryOptions) (*nomadApi.Job, *nomadApi.QueryMeta, error)
	Register(job *nomadApi.Job, q *nomadApi.WriteOptions) (*nomadApi.JobRegisterResponse, *nomadApi.WriteMeta, error)
}

// NomadDeployments represents the interface for interacting
// with nomads over the deployments end-point
type NomadDeployments interface {
	Info(deploymentID string, q *nomadApi.QueryOptions) (*nomadApi.Deployment, *nomadApi.QueryMeta, error)
}

// NomadEvaluations represents the interface for interacting
// with nomads over the evaluations end-point
type NomadEvaluations interface {
	Info(evalID string, q *nomadApi.QueryOptions) (*nomadApi.Evaluation, *nomadApi.QueryMeta, error)
}
