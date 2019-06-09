package nomad

import (
	nomadApi "github.com/hashicorp/nomad/api"
)

// Jobs represents the interface for interacting
// with nomads over the jobs end-point
type Jobs interface {
	Info(jobID string, q *nomadApi.QueryOptions) (*nomadApi.Job, *nomadApi.QueryMeta, error)
	Register(job *nomadApi.Job, q *nomadApi.WriteOptions) (*nomadApi.JobRegisterResponse, *nomadApi.WriteMeta, error)
}

// Deployments represents the interface for interacting
// with nomads over the deployments end-point
type Deployments interface {
	Info(deploymentID string, q *nomadApi.QueryOptions) (*nomadApi.Deployment, *nomadApi.QueryMeta, error)
}

// Evaluations represents the interface for interacting
// with nomads over the evaluations end-point
type Evaluations interface {
	Info(evalID string, q *nomadApi.QueryOptions) (*nomadApi.Evaluation, *nomadApi.QueryMeta, error)
}
