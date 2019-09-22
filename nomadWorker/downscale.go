package nomadWorker

import (
	"fmt"

	nomadApi "github.com/hashicorp/nomad/api"
)

type candidate struct {
	// nodeID is the nomad node ID
	nodeID     string
	datacenter string
	// instanceID is the aws instance id
	instanceID string
	ipAddress  string
}

func (c *Connector) downscale(datacenter string, desiredCount uint) error {

	// 1. Select a candidate for downscaling -> returns [needs node id]
	candidate, err := selectCandidate(c.nodesIF, datacenter)
	if err != nil {
		return err
	}

	c.log.Info().Str("NodeID", candidate.nodeID).Msgf("1. [Select] Selected node '%s' (%s, %s) as candidate for downscaling.", candidate.nodeID, candidate.ipAddress, candidate.instanceID)

	// 2. Drain the node [needs node id]
	c.log.Info().Str("NodeID", candidate.nodeID).Msgf("2. [Drain] Draining node '%s' (%s, %s) ... ", candidate.nodeID, candidate.ipAddress, candidate.instanceID)
	nodeModifyIndex, err := drainNode(c.nodesIF, candidate.nodeID, c.nodeDrainDeadline)
	if err != nil {
		return err
	}
	monitorDrainNode(c.nodesIF, candidate.nodeID, nodeModifyIndex, c.log)
	c.log.Info().Str("NodeID", candidate.nodeID).Msgf("2. [Drain] Draining node '%s' (%s, %s) ... done", candidate.nodeID, candidate.ipAddress, candidate.instanceID)

	// 3. Terminate the node using the AWS ASG [needs instance id]

	if err := setEligibility(c.nodesIF, candidate.nodeID, true); err != nil {
		return err
	}
	return nil
}

func setEligibility(nodesIF Nodes, nodeID string, eligible bool) error {
	_, err := nodesIF.ToggleEligibility(nodeID, eligible, nil)
	return err
}

func selectCandidate(nodesIF Nodes, datacenter string) (*candidate, error) {

	nodeListStub, _, err := nodesIF.List(nil)
	if err != nil {
		return nil, err
	}

	// filter out the nodes for this datacenter that are not draining already and are ready
	nodes := make([]*nomadApi.NodeListStub, 0)
	for _, node := range nodeListStub {
		if !node.Drain && node.Datacenter == datacenter && node.Status == nomadApi.NodeStatusReady {
			nodes = append(nodes, node)
		}
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("No node found in datacenter '%s' that is not already draining", datacenter)
	}

	// now select the best node
	// TODO: select the node with least running allocations
	// Hint: https://www.nomadproject.io/api/nodes.html#list-node-allocations
	// HACK: Just take the first node for now
	node := nodes[0]

	return &candidate{
		ipAddress:  node.Address,
		instanceID: node.Name,
		nodeID:     node.ID,
		datacenter: node.Datacenter,
	}, nil
}
