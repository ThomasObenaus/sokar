package nomadWorker

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
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
