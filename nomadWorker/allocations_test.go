package nomadWorker

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/nomadWorker"
)

func Test_GetNumAllocationsInStatus(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)

	nodeID := "nodeID"
	// error
	nodesIF.EXPECT().Allocations(nodeID, nil).Return(nil, nil, fmt.Errorf("ERRR"))
	num, err := getNumAllocationsInStatus(nodesIF, nodeID, "running")
	assert.Error(t, err)
	assert.Equal(t, uint(0), num)

	// success
	allocations := make([]*nomadApi.Allocation, 0)
	statusRunning := "running"
	jobRunning := nomadApi.Job{Status: &statusRunning}
	allocation := &nomadApi.Allocation{Job: &jobRunning}
	allocations = append(allocations, allocation)
	statusStopped := "stopped"
	jobStopped := nomadApi.Job{Status: &statusStopped}
	allocation = &nomadApi.Allocation{Job: &jobStopped}
	allocations = append(allocations, allocation)
	nodesIF.EXPECT().Allocations(nodeID, nil).Return(allocations, nil, nil)
	num, err = getNumAllocationsInStatus(nodesIF, nodeID, "running")
	assert.NoError(t, err)
	assert.Equal(t, uint(1), num)
}
