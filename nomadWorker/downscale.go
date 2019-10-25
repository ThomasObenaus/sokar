package nomadWorker

import (
	"fmt"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/aws"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

type candidate struct {
	// nodeID is the nomad node ID
	nodeID     string
	datacenter string
	// instanceID is the aws instance id
	instanceID            string
	ipAddress             string
	numRunningAllocations uint
	cpu                   int
	diskMB                int
	memoryMB              int
}

func (c *Connector) downscale(datacenter string, desiredCount uint) error {

	// 1. Select a candidate for downscaling -> returns [needs node id]
	candidate, err := selectCandidate(c.nodesIF, datacenter, c.log)
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
	c.log.Info().Str("NodeID", candidate.nodeID).Msgf("3. [Terminate] Terminate node '%s' (%s, %s) ... ", candidate.nodeID, candidate.ipAddress, candidate.instanceID)
	sess, err := c.createSession()
	if err != nil {
		return err
	}
	autoScalingIF := c.autoScalingFactory.CreateAutoScaling(sess)
	autoscalingGroupName, activityID, err := aws.TerminateInstanceInAsg(autoScalingIF, candidate.instanceID)
	if err != nil {
		return err
	}

	// wait until the instance is scaled down
	if iter, err := aws.MonitorInstanceScaling(autoScalingIF, autoscalingGroupName, activityID, c.monitorInstanceTimeout, c.log); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Monitor instance scaling failed after %d iterations.", iter))
	}
	c.log.Info().Str("NodeID", candidate.nodeID).Msgf("3. [Terminate] Terminate node '%s' (%s, %s) ... done", candidate.nodeID, candidate.ipAddress, candidate.instanceID)
	return nil
}

func setEligibility(nodesIF Nodes, nodeID string, eligible bool) error {
	_, err := nodesIF.ToggleEligibility(nodeID, eligible, nil)
	return errors.WithMessage(err, "Failed toggling node eligibility")
}

func selectCandidate(nodesIF Nodes, datacenter string, log zerolog.Logger) (*candidate, error) {

	nodeListStub, _, err := nodesIF.List(nil)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed listing nomad nodes")
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

	// now select the best node based on the least running allocations
	bestCandidate := candidate{
		numRunningAllocations: MaxUint,
		cpu:                   MaxInt,
		memoryMB:              MaxInt,
		diskMB:                MaxInt,
	}

	for _, node := range nodes {
		allocInfo, err := getNumAllocationsInStatus(nodesIF, node.ID, nomadApi.AllocClientStatusRunning)
		if err != nil {
			log.Error().Err(err).Msg("Unable to obtain the nodes allocations. Ignore node.")
			continue
		}

		lessAllocs := allocInfo.numAllocations < bestCandidate.numRunningAllocations
		sameAllocs := allocInfo.numAllocations == bestCandidate.numRunningAllocations
		sameMem := allocInfo.memoryMB == bestCandidate.memoryMB
		sameAllocsButLessMem := sameAllocs && (allocInfo.memoryMB < bestCandidate.memoryMB)
		sameAllocsAndMemButCPU := sameAllocs && sameMem && (allocInfo.cpu < bestCandidate.cpu)

		if lessAllocs || sameAllocsButLessMem || sameAllocsAndMemButCPU {
			bestCandidate.nodeID = node.ID
			bestCandidate.numRunningAllocations = allocInfo.numAllocations
			bestCandidate.cpu = allocInfo.cpu
			bestCandidate.memoryMB = allocInfo.memoryMB
			bestCandidate.diskMB = allocInfo.diskMB
			bestCandidate.datacenter = node.Datacenter
			bestCandidate.instanceID = node.Name
			bestCandidate.ipAddress = node.Address
			log.Info().Msgf("New best candidate (nodeID=%s) with %d running allocations found (cpu=%d,diskMB=%d,memMB=%d).", bestCandidate.nodeID, bestCandidate.numRunningAllocations, bestCandidate.cpu, bestCandidate.diskMB, bestCandidate.memoryMB)
		}
	}

	return &bestCandidate, nil
}
