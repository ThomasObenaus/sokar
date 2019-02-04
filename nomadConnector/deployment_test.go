package nomadConnector

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/nomadConnector"
)

func TestGetDeploymentID_NoIF(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	conn := connectorImpl{}

	deplID, err := conn.getDeploymentID("ABCDEF", time.Millisecond*600)
	assert.Error(t, err)
	assert.Empty(t, deplID)
}

func TestGetDeploymentID_Success(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evalIF := mock_nomadConnector.NewMockNomadEvaluations(mockCtrl)
	conn := connectorImpl{
		evalIF: evalIF,
	}

	// success
	deplIDWanted := "DEPL1234"
	eval := nomadApi.Evaluation{DeploymentID: deplIDWanted}
	evalID := "ABCDEFG"
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, nil)

	deplID, err := conn.getDeploymentID(evalID, time.Millisecond*600)
	assert.NoError(t, err)
	assert.Equal(t, deplIDWanted, deplID)
}

func TestGetDeploymentID_Timeout(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evalIF := mock_nomadConnector.NewMockNomadEvaluations(mockCtrl)
	conn := connectorImpl{
		evalIF: evalIF,
	}

	// timeout
	evalID := "ABCDEFG"
	eval := nomadApi.Evaluation{}
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, nil).AnyTimes()

	deplID, err := conn.getDeploymentID(evalID, time.Millisecond*600)
	assert.Error(t, err)
	assert.Empty(t, deplID)
}

func TestGetDeploymentID_InternalErr(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evalIF := mock_nomadConnector.NewMockNomadEvaluations(mockCtrl)
	conn := connectorImpl{
		evalIF: evalIF,
	}

	// timeout
	evalID := "ABCDEFG"
	eval := nomadApi.Evaluation{}
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, fmt.Errorf("Internal error")).AnyTimes()

	deplID, err := conn.getDeploymentID(evalID, time.Millisecond*600)
	assert.Error(t, err)
	assert.Empty(t, deplID)
}
