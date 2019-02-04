package nomadConnector

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/nomadConnector"
)

func TestGetJobInfo(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// interface missing test
	conn := connectorImpl{}

	jobInfo, err := conn.getJobInfo("unknown")
	assert.Error(t, err)
	assert.Nil(t, jobInfo)

	jobsIF := mock_nomadConnector.NewMockNomadJobs(mockCtrl)
	conn = connectorImpl{
		jobsIF: jobsIF,
	}

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

func TestGetJobCount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jobsIF := mock_nomadConnector.NewMockNomadJobs(mockCtrl)
	conn := connectorImpl{
		jobsIF: jobsIF,
	}

	// count 0
	job := &nomadApi.Job{}
	jobsIF.EXPECT().Info("test", &nomadApi.QueryOptions{AllowStale: true}).Return(job, nil, nil)

	count, err := conn.GetJobCount("test")
	assert.NoError(t, err)
	assert.Equal(t, uint(0), count)

	// count 10
	count10 := 10
	count5 := 5
	job = &nomadApi.Job{
		TaskGroups: []*nomadApi.TaskGroup{{Count: &count10}, {Count: &count5}},
	}
	jobsIF.EXPECT().Info("test", &nomadApi.QueryOptions{AllowStale: true}).Return(job, nil, nil)

	count, err = conn.GetJobCount("test")
	assert.NoError(t, err)
	assert.Equal(t, uint(10), count)
}
