package nomadWorker

import (
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

func monitorDrainNode(nodesIF Nodes, nodeID string, nodeModifyIndex uint64, logger zerolog.Logger) {

	logger.Info().Str("NodeID", nodeID).Msgf("Monitoring node draining (node=%s) ... ", nodeID)

	deadline := time.Now().Add(time.Second * 60)
	ctx := monitoringCtx{
		deadline: deadline,
		doneChan: make(chan struct{}),
	}

	events := nodesIF.MonitorDrain(ctx, nodeID, nodeModifyIndex, false)
	for ev := range events {
		logger.Info().Str("NodeID", nodeID).Msg(ev.String())
	}
	logger.Info().Str("NodeID", nodeID).Msgf("Monitoring node draining (node=%s) ... done", nodeID)
}

type monitoringCtx struct {
	doneChan <-chan struct{}
	deadline time.Time
}

func (ctx monitoringCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.deadline, false
}
func (ctx monitoringCtx) Done() <-chan struct{} {
	return ctx.doneChan
}
func (ctx monitoringCtx) Err() error {
	return nil
}
func (ctx monitoringCtx) Value(key interface{}) interface{} {
	return nil
}
