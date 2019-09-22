package nomadWorker

import (
	nomadApi "github.com/hashicorp/nomad/api"
)

// Nodes represents the minimal interface used to gather information about nomad nodes
type Nodes interface {
	List(q *nomadApi.QueryOptions) ([]*nomadApi.NodeListStub, *nomadApi.QueryMeta, error)
	ToggleEligibility(nodeID string, eligible bool, q *nomadApi.WriteOptions) (*nomadApi.NodeEligibilityUpdateResponse, error)
	UpdateDrain(nodeID string, spec *nomadApi.DrainSpec, markEligible bool, q *nomadApi.WriteOptions) (*nomadApi.NodeDrainUpdateResponse, error)
}
