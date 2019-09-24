package nomadWorker

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_aws "github.com/thomasobenaus/sokar/test/aws"
	"github.com/thomasobenaus/sokar/test/nomadWorker"
)

func TestSelectCandidateForDownscaling_Errors(t *testing.T) {

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)
	datacenter := "dcXYZ"
	// no nodes
	nodes := make([]*nomadApi.NodeListStub, 0)
	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	candidate, err := selectCandidate(nodesIF, datacenter, logger)
	assert.Nil(t, candidate)
	assert.Error(t, err)

	// no nodes in datacenter
	nodes = make([]*nomadApi.NodeListStub, 0)
	node := nomadApi.NodeListStub{Datacenter: "other_dc"}
	nodes = append(nodes, &node)
	qmeta = nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	candidate, err = selectCandidate(nodesIF, datacenter, logger)
	assert.Nil(t, candidate)
	assert.Error(t, err)

	// no nodes in datacenter that are not draining
	nodes = make([]*nomadApi.NodeListStub, 0)
	node = nomadApi.NodeListStub{Datacenter: datacenter, Drain: true}
	nodes = append(nodes, &node)
	qmeta = nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	candidate, err = selectCandidate(nodesIF, datacenter, logger)
	assert.Nil(t, candidate)
	assert.Error(t, err)

	// valid nodes available but down
	nodes = make([]*nomadApi.NodeListStub, 0)
	node = nomadApi.NodeListStub{Datacenter: datacenter, Drain: false, Name: "node1", ID: "1234", Address: "192.1680.0.1", Status: "down"}
	nodes = append(nodes, &node)
	qmeta = nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	candidate, err = selectCandidate(nodesIF, datacenter, logger)
	assert.Nil(t, candidate)
	assert.Error(t, err)
}

func TestSelectCandidateForDownscaling_Success(t *testing.T) {

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)
	datacenter := "dcXYZ"

	// valid nodes available
	nodeID := "nodeID"
	nodes := make([]*nomadApi.NodeListStub, 0)
	node := nomadApi.NodeListStub{Datacenter: datacenter, Drain: false, Name: "node1", ID: nodeID, Address: "192.1680.0.1", Status: "ready"}
	nodes = append(nodes, &node)
	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)

	allocations := make([]*nomadApi.Allocation, 0)
	statusRunning := "running"
	jobRunning := nomadApi.Job{Status: &statusRunning}
	allocation := &nomadApi.Allocation{Job: &jobRunning}
	allocations = append(allocations, allocation)
	nodesIF.EXPECT().Allocations(nodeID, nil).Return(allocations, nil, nil)
	candidate, err := selectCandidate(nodesIF, datacenter, logger)
	assert.NotNil(t, candidate)
	assert.Equal(t, nodeID, candidate.nodeID)
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

func Test_Downscale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgFactory := mock_aws.NewMockAutoScalingFactory(mockCtrl)
	nodesIF := mock_nomadWorker.NewMockNodes(mockCtrl)
	asgIF := mock_aws.NewMockAutoScaling(mockCtrl)

	cfg := Config{AWSProfile: "xyz", NomadServerAddress: "http://nomad.io"}
	connector, err := cfg.New()
	require.NotNil(t, connector)
	require.NoError(t, err)

	connector.autoScalingFactory = asgFactory
	connector.nodesIF = nodesIF

	instanceID := "1234"
	datacenter := "private-services"
	desiredCount := uint(3)
	nodes := make([]*nomadApi.NodeListStub, 0)
	node := nomadApi.NodeListStub{Datacenter: datacenter, Drain: false, Name: instanceID, ID: "1234", Address: "192.1680.0.1", Status: "ready"}
	nodes = append(nodes, &node)
	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	nodesIF.EXPECT().List(gomock.Any()).Return(nodes, &qmeta, nil)
	nodeID := "1234"
	nodeModifyIndex := uint64(1234)
	nodeDrainResp := nomadApi.NodeDrainUpdateResponse{NodeModifyIndex: nodeModifyIndex}
	nodesIF.EXPECT().UpdateDrain(nodeID, gomock.Any(), false, nil).Return(&nodeDrainResp, nil)

	evChan := make(chan *nomadApi.MonitorMessage)
	msg := nomadApi.MonitorMessage{}

	go func() {
		evChan <- &msg
		close(evChan)
	}()

	nodesIF.EXPECT().MonitorDrain(gomock.Any(), nodeID, nodeModifyIndex, false).Return(evChan)
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)

	shouldDecDesiredCapa := true
	input := autoscaling.TerminateInstanceInAutoScalingGroupInput{InstanceId: &instanceID, ShouldDecrementDesiredCapacity: &shouldDecDesiredCapa}
	req := request.Request{}
	activityID := "ActivityId"
	asgName := "AsgName"

	progress := int64(100)
	statusCode := "InProgress"
	activity := autoscaling.Activity{ActivityId: &activityID, AutoScalingGroupName: &asgName, Progress: &progress, StatusCode: &statusCode}
	output := autoscaling.TerminateInstanceInAutoScalingGroupOutput{Activity: &activity}
	asgIF.EXPECT().TerminateInstanceInAutoScalingGroupRequest(&input).Return(&req, &output)

	activities := make([]*autoscaling.Activity, 0)
	activities = append(activities, &activity)
	outputDescribeActivities := autoscaling.DescribeScalingActivitiesOutput{Activities: activities}
	asgIF.EXPECT().DescribeScalingActivitiesRequest(gomock.Any()).Return(&req, &outputDescribeActivities)

	allocations := make([]*nomadApi.Allocation, 0)
	statusRunning := "running"
	jobRunning := nomadApi.Job{Status: &statusRunning}
	allocation := &nomadApi.Allocation{Job: &jobRunning}
	allocations = append(allocations, allocation)
	nodesIF.EXPECT().Allocations(nodeID, nil).Return(allocations, nil, nil)
	err = connector.downscale(datacenter, desiredCount)
	assert.NoError(t, err)

}
