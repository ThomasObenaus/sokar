package nomadWorker

import (
	"testing"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/nomadWorker"
)

func TestSelectCandidateForDownscaling_Errors(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)
	datacenter := "dcXYZ"
	// no nodes
	nodes := make([]*nomadApi.NodeListStub, 0)
	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	candidate, err := selectCandidate(nodesIF, datacenter)
	assert.Nil(t, candidate)
	assert.Error(t, err)

	// no nodes in datacenter
	nodes = make([]*nomadApi.NodeListStub, 0)
	node := nomadApi.NodeListStub{Datacenter: "other_dc"}
	nodes = append(nodes, &node)
	qmeta = nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	candidate, err = selectCandidate(nodesIF, datacenter)
	assert.Nil(t, candidate)
	assert.Error(t, err)

	// no nodes in datacenter that are not draining
	nodes = make([]*nomadApi.NodeListStub, 0)
	node = nomadApi.NodeListStub{Datacenter: datacenter, Drain: true}
	nodes = append(nodes, &node)
	qmeta = nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	candidate, err = selectCandidate(nodesIF, datacenter)
	assert.Nil(t, candidate)
	assert.Error(t, err)

	// valid nodes available but down
	nodes = make([]*nomadApi.NodeListStub, 0)
	node = nomadApi.NodeListStub{Datacenter: datacenter, Drain: false, Name: "node1", ID: "1234", Address: "192.1680.0.1", Status: "down"}
	nodes = append(nodes, &node)
	qmeta = nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	candidate, err = selectCandidate(nodesIF, datacenter)
	assert.Nil(t, candidate)
	assert.Error(t, err)
}

func TestSelectCandidateForDownscaling_Success(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)
	datacenter := "dcXYZ"

	// valid nodes available
	nodes := make([]*nomadApi.NodeListStub, 0)
	node := nomadApi.NodeListStub{Datacenter: datacenter, Drain: false, Name: "node1", ID: "1234", Address: "192.1680.0.1", Status: "ready"}
	nodes = append(nodes, &node)
	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	candidate, err := selectCandidate(nodesIF, datacenter)
	assert.NotNil(t, candidate)
	assert.Equal(t, "1234", candidate.nodeID)
	assert.NoError(t, err)
}

func TestSetEligibility(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)

	nodeID := "1234"
	nodesIF.EXPECT().ToggleEligibility(nodeID, true, nil).Return(nil, nil)
	err := setEligibility(nodesIF, nodeID, true)
	assert.NoError(t, err)
}
