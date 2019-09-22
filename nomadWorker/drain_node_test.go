package nomadWorker

import (
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/nomadWorker"
)

func TestDrainNode(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)

	nodeDrainResp := nomadApi.NodeDrainUpdateResponse{NodeModifyIndex: 1234}

	nodeID := "1234"
	nodesIF.EXPECT().UpdateDrain(nodeID, gomock.Any(), false, nil).Return(&nodeDrainResp, nil)
	idx, err := drainNode(nodesIF, nodeID, time.Second*20)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1234), idx)
}
func TestMonitorDrainNode(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	nodeID := "1234"
	nodeModifyIndex := uint64(1234)
	evChan := make(chan *nomadApi.MonitorMessage)
	msg := nomadApi.MonitorMessage{}

	go func() {
		evChan <- &msg
		close(evChan)
	}()

	nodesIF.EXPECT().MonitorDrain(gomock.Any(), nodeID, nodeModifyIndex, false).Return(evChan)
	numEvents := monitorDrainNode(nodesIF, nodeID, nodeModifyIndex, logger)
	assert.Equal(t, uint(1), numEvents)
}
