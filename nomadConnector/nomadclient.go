package nomadConnector

import (
	nomadApi "github.com/hashicorp/nomad/api"
)

// NomadClient represents an interface providing the methods
// to interact with nomad
type NomadClient interface {
	Deployments() *nomadApi.Deployments
	Evaluations() *nomadApi.Evaluations
	Jobs() *nomadApi.Jobs
}
