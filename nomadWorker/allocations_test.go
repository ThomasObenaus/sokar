package nomadWorker

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_nomadWorker "github.com/thomasobenaus/sokar/test/nomadWorker"
)

func Test_GetNumAllocationsInStatus(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)

	nodeID := "nodeID"
	// error
	nodesIF.EXPECT().Allocations(nodeID, nil).Return(nil, nil, fmt.Errorf("ERRR"))
	allocInfo, err := getNumAllocationsInStatus(nodesIF, nodeID, "running")
	assert.Error(t, err)
	assert.Nil(t, allocInfo)

	// success
	resourceValue := 10
	resources := nomadApi.Resources{CPU: &resourceValue, MemoryMB: &resourceValue, DiskMB: &resourceValue}
	allocations := make([]*nomadApi.Allocation, 0)
	statusRunning := "running"
	jobRunning := nomadApi.Job{Status: &statusRunning}
	allocation := &nomadApi.Allocation{Job: &jobRunning, Resources: &resources}
	allocations = append(allocations, allocation)
	statusStopped := "stopped"
	jobStopped := nomadApi.Job{Status: &statusStopped}
	allocation = &nomadApi.Allocation{Job: &jobStopped, Resources: &resources}
	allocations = append(allocations, allocation)
	nodesIF.EXPECT().Allocations(nodeID, nil).Return(allocations, nil, nil)
	allocInfo, err = getNumAllocationsInStatus(nodesIF, nodeID, "running")
	require.NotNil(t, allocInfo)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), allocInfo.numAllocations)
	assert.Equal(t, 10, allocInfo.cpu)
	assert.Equal(t, 10, allocInfo.diskMB)
	assert.Equal(t, 10, allocInfo.memoryMB)
}
