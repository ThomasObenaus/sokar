package nomadWorker

import (
	"context"
	"time"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
)

func drainNode(nodesIF Nodes, nodeID string, deadline time.Duration) (nodeModifyIndex uint64, err error) {

	drainSpec := nomadApi.DrainSpec{
		Deadline:         deadline,
		IgnoreSystemJobs: false,
	}

	resp, err := nodesIF.UpdateDrain(nodeID, &drainSpec, false, nil)
	return resp.NodeModifyIndex, err
}

func monitorDrainNode(nodesIF Nodes, nodeID string, nodeModifyIndex uint64, timeout time.Duration, logger zerolog.Logger) uint {

	logger.Info().Str("NodeID", nodeID).Msgf("Monitoring node draining (node=%s, timeout=%s) ... ", nodeID, timeout.String())

	var numEvents uint
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// FIXME: Find out if and when we need to call the cancel function from the outside to close the context
	_ = cancel

	// create and obtain the monitoring channel and then wait until it is closed
	events := nodesIF.MonitorDrain(ctx, nodeID, nodeModifyIndex, false)
	for ev := range events {
		if ev != nil {
			logger.Info().Str("NodeID", nodeID).Msg(ev.String())
			numEvents++
		}
	}
	logger.Info().Str("NodeID", nodeID).Msgf("Monitoring node draining (node=%s, timeout=%s, #events=%d) ... done", nodeID, timeout.String(), numEvents)
	return numEvents
}
