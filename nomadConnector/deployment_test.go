package nomadConnector

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/nomadConnector"
)

func TestDeployment(t *testing.T) {

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
