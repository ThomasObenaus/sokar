package nomadConnector

import (
	nomadApi "github.com/hashicorp/nomad/api"
)

// NomadClient represents an interface providing the methods
// to interact with nomad
type NomadClient interface {
	Deployments() *nomadApi.Deployments
	Evaluations() *nomadApi.Evaluations
}

type NomadJobs interface {
	Info(jobID string, q *nomadApi.QueryOptions) (*nomadApi.Job, *nomadApi.QueryMeta, error)
	Register(job *nomadApi.Job, q *nomadApi.WriteOptions) (*nomadApi.JobRegisterResponse, *nomadApi.WriteMeta, error)
}
