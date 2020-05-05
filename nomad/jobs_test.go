package nomad

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"

	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/nomad/structs"
	mock_nomad "github.com/thomasobenaus/sokar/test/nomad"
)

func minimalConnectorImpl() Connector {
	conn := Connector{
		deploymentTimeOut: time.Second * 20,
		evaluationTimeOut: time.Second * 10,
	}
	return conn
}

func TestGetJobInfo(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// interface missing test
	conn := minimalConnectorImpl()

	jobInfo, err := conn.getJobInfo("unknown")
	assert.Error(t, err)
	assert.Nil(t, jobInfo)

	jobsIF := mock_nomad.NewMockJobs(mockCtrl)
	conn.jobsIF = jobsIF

	// job not found test
	jobsIF.EXPECT().Info("unknown", &nomadApi.QueryOptions{AllowStale: true}).Return(nil, nil, fmt.Errorf("Job not found"))

	jobInfo, err = conn.getJobInfo("unknown")
	assert.Error(t, err)
	assert.Nil(t, jobInfo)

	// job found test
	job := &nomadApi.Job{}
	jobsIF.EXPECT().Info("test", &nomadApi.QueryOptions{AllowStale: true}).Return(job, nil, nil)

	jobInfo, err = conn.getJobInfo("test")
	assert.NoError(t, err)
	assert.NotNil(t, jobInfo)
}

func TestAdjustScalingObjectCount_Success(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jobsIF := mock_nomad.NewMockJobs(mockCtrl)
	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	deplIF := mock_nomad.NewMockDeployments(mockCtrl)

	conn := minimalConnectorImpl()
	conn.evalIF = evalIF
	conn.deploymentIF = deplIF
	conn.jobsIF = jobsIF

	// GetJobInfo
	count10 := 10
	count5 := 5
	nameA := "grpA"
	nameB := "grpB"
	job := &nomadApi.Job{
		TaskGroups: []*nomadApi.TaskGroup{{Name: &nameA, Count: &count10}, {Name: &nameB, Count: &count5}},
	}
	jobsIF.EXPECT().Info("test", &nomadApi.QueryOptions{AllowStale: true}).Return(job, nil, nil)

	// Register Deployment
	jobRegisterResponse := nomadApi.JobRegisterResponse{EvalID: "ABCDEFG"}
	jobsIF.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&jobRegisterResponse, nil, nil)

	// Obtain DeplyomentID
	deplID := "DEPL1234"
	eval := nomadApi.Evaluation{DeploymentID: deplID}
	evalID := "ABCDEFG"
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, nil)

	// Wait for deployment confirmation
	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	depl := nomadApi.Deployment{Status: structs.DeploymentStatusSuccessful}
	deplIF.EXPECT().Info(deplID, gomock.Any()).Return(&depl, &qmeta, nil)

	err := conn.AdjustScalingObjectCount("test", 2, 10, 4, 5)
	assert.NoError(t, err)
}

func TestAdjustScalingObjectCount_InternalError(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jobsIF := mock_nomad.NewMockJobs(mockCtrl)
	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	deplIF := mock_nomad.NewMockDeployments(mockCtrl)

	conn := minimalConnectorImpl()
	conn.evalIF = evalIF
	conn.deploymentIF = deplIF
	conn.jobsIF = jobsIF

	// GetJobInfo
	count10 := 10
	count5 := 5
	nameA := "grpA"
	nameB := "grpB"
	job := &nomadApi.Job{
		TaskGroups: []*nomadApi.TaskGroup{{Name: &nameA, Count: &count10}, {Name: &nameB, Count: &count5}},
	}
	jobsIF.EXPECT().Info("test", &nomadApi.QueryOptions{AllowStale: true}).Return(job, nil, nil)

	// Register Deployment
	jobsIF.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, nil, fmt.Errorf("Internal error"))

	err := conn.AdjustScalingObjectCount("test", 2, 10, 4, 5)
	assert.Error(t, err)
}

func TestAdjustScalingObjectCount_DeploymentError(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jobsIF := mock_nomad.NewMockJobs(mockCtrl)
	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	deplIF := mock_nomad.NewMockDeployments(mockCtrl)

	conn := minimalConnectorImpl()
	conn.evalIF = evalIF
	conn.deploymentIF = deplIF
	conn.jobsIF = jobsIF

	// GetJobInfo
	count10 := 10
	count5 := 5
	nameA := "grpA"
	nameB := "grpB"
	job := &nomadApi.Job{
		TaskGroups: []*nomadApi.TaskGroup{{Name: &nameA, Count: &count10}, {Name: &nameB, Count: &count5}},
	}
	jobsIF.EXPECT().Info("test", &nomadApi.QueryOptions{AllowStale: true}).Return(job, nil, nil)

	// Register Deployment
	jobRegisterResponse := nomadApi.JobRegisterResponse{EvalID: "ABCDEFG"}
	jobsIF.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&jobRegisterResponse, nil, nil)

	// Obtain DeplyomentID
	deplID := "DEPL1234"
	eval := nomadApi.Evaluation{DeploymentID: deplID}
	evalID := "ABCDEFG"
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, nil)

	// Wait for deployment confirmation
	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	depl := nomadApi.Deployment{Status: structs.DeploymentStatusCancelled}
	deplIF.EXPECT().Info(deplID, gomock.Any()).Return(&depl, &qmeta, nil)

	err := conn.AdjustScalingObjectCount("test", 2, 10, 4, 5)
	assert.Error(t, err)
}

func TestGetScalingObjectCount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jobsIF := mock_nomad.NewMockJobs(mockCtrl)
	conn := minimalConnectorImpl()
	conn.jobsIF = jobsIF

	// count 0
	job := &nomadApi.Job{}
	jobsIF.EXPECT().Info("test", &nomadApi.QueryOptions{AllowStale: true}).Return(job, nil, nil)

	count, err := conn.GetScalingObjectCount("test")
	assert.NoError(t, err)
	assert.Equal(t, uint(0), count)

	// count 10
	count10 := 10
	count5 := 5
	job = &nomadApi.Job{
		TaskGroups: []*nomadApi.TaskGroup{{Count: &count10}, {Count: &count5}},
	}
	jobsIF.EXPECT().Info("test", &nomadApi.QueryOptions{AllowStale: true}).Return(job, nil, nil)

	count, err = conn.GetScalingObjectCount("test")
	assert.NoError(t, err)
	assert.Equal(t, uint(10), count)
}

func TestIsScalingObjectDead(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jobsIF := mock_nomad.NewMockJobs(mockCtrl)
	conn := minimalConnectorImpl()
	conn.jobsIF = jobsIF

	// fail
	job := &nomadApi.Job{}
	jobsIF.EXPECT().Info("test", gomock.Any()).Return(job, nil, nil)

	dead, err := conn.IsScalingObjectDead("test")
	assert.Error(t, err)
	assert.Equal(t, false, dead)

	// success
	status := structs.JobStatusDead
	job = &nomadApi.Job{Status: &status}
	jobsIF.EXPECT().Info("test", gomock.Any()).Return(job, nil, nil)

	dead, err = conn.IsScalingObjectDead("test")
	assert.NoError(t, err)
	assert.Equal(t, true, dead)
}
