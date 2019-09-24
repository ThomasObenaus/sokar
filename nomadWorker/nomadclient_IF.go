package nomadWorker

import (
	"context"

	nomadApi "github.com/hashicorp/nomad/api"
)

// Nodes represents the minimal interface used to gather information about nomad nodes
type Nodes interface {
	List(q *nomadApi.QueryOptions) ([]*nomadApi.NodeListStub, *nomadApi.QueryMeta, error)
	ToggleEligibility(nodeID string, eligible bool, q *nomadApi.WriteOptions) (*nomadApi.NodeEligibilityUpdateResponse, error)
	UpdateDrain(nodeID string, spec *nomadApi.DrainSpec, markEligible bool, q *nomadApi.WriteOptions) (*nomadApi.NodeDrainUpdateResponse, error)
	MonitorDrain(ctx context.Context, nodeID string, index uint64, ignoreSys bool) <-chan *nomadApi.MonitorMessage
	Allocations(nodeID string, q *nomadApi.QueryOptions) ([]*nomadApi.Allocation, *nomadApi.QueryMeta, error)
}
