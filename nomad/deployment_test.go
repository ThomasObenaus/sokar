package nomad

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/thomasobenaus/sokar/nomad/structs"
	mock_nomad "github.com/thomasobenaus/sokar/test/nomad"

	"github.com/stretchr/testify/assert"
)

func TestGetDeploymentID_NoIF(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	conn := Connector{}

	deplID, err := conn.getDeploymentID("ABCDEF", time.Millisecond*600)
	assert.Error(t, err)
	assert.Empty(t, deplID)
}

func TestWaitForDeploymentConfirmation_Success(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	deplIF := mock_nomad.NewMockDeployments(mockCtrl)
	conn := minimalConnectorImpl()
	conn.evalIF = evalIF
	conn.deploymentIF = deplIF

	deplID := "DEPL1234"
	eval := nomadApi.Evaluation{DeploymentID: deplID}
	evalID := "ABCDEFG"
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, nil)

	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	depl := nomadApi.Deployment{Status: structs.DeploymentStatusSuccessful}
	deplIF.EXPECT().Info(deplID, gomock.Any()).Return(&depl, &qmeta, nil)

	err := conn.waitForDeploymentConfirmation(evalID, time.Millisecond*600)
	assert.NoError(t, err)
}

func TestWaitForDeploymentConfirmation_Timeout(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	deplIF := mock_nomad.NewMockDeployments(mockCtrl)
	conn := minimalConnectorImpl()
	conn.evalIF = evalIF
	conn.deploymentIF = deplIF

	deplID := "DEPL1234"
	eval := nomadApi.Evaluation{DeploymentID: deplID}
	evalID := "ABCDEFG"
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, nil)

	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	depl := nomadApi.Deployment{Status: structs.DeploymentStatusRunning}
	deplIF.EXPECT().Info(deplID, gomock.Any()).Return(&depl, &qmeta, nil)

	err := conn.waitForDeploymentConfirmation(evalID, time.Millisecond*600)
	assert.Error(t, err)
}

func TestWaitForDeploymentConfirmation_Failed(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	deplIF := mock_nomad.NewMockDeployments(mockCtrl)
	conn := minimalConnectorImpl()
	conn.evalIF = evalIF
	conn.deploymentIF = deplIF

	deplID := "DEPL1234"
	eval := nomadApi.Evaluation{DeploymentID: deplID}
	evalID := "ABCDEFG"
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, nil)

	qmeta := nomadApi.QueryMeta{LastIndex: 1000}
	depl := nomadApi.Deployment{Status: structs.DeploymentStatusCancelled}
	deplIF.EXPECT().Info(deplID, gomock.Any()).Return(&depl, &qmeta, nil)

	err := conn.waitForDeploymentConfirmation(evalID, time.Millisecond*600)
	assert.Error(t, err)
}

func TestWaitForDeploymentConfirmation_Nil(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	deplIF := mock_nomad.NewMockDeployments(mockCtrl)
	conn := minimalConnectorImpl()
	conn.evalIF = evalIF
	conn.deploymentIF = deplIF

	deplID := "DEPL1234"
	eval := nomadApi.Evaluation{DeploymentID: deplID}
	evalID := "ABCDEFG"
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, nil)

	depl := nomadApi.Deployment{Status: structs.DeploymentStatusCancelled}
	deplIF.EXPECT().Info(deplID, gomock.Any()).Return(&depl, nil, nil)

	err := conn.waitForDeploymentConfirmation(evalID, time.Millisecond*600)
	assert.Error(t, err)
}

func TestGetDeploymentID_Success(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	conn := minimalConnectorImpl()
	conn.evalIF = evalIF

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

	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	conn := minimalConnectorImpl()
	conn.evalIF = evalIF

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

	evalIF := mock_nomad.NewMockEvaluations(mockCtrl)
	conn := minimalConnectorImpl()
	conn.evalIF = evalIF

	// timeout
	evalID := "ABCDEFG"
	eval := nomadApi.Evaluation{}
	evalIF.EXPECT().Info(evalID, nil).Return(&eval, nil, fmt.Errorf("Internal error")).AnyTimes()

	deplID, err := conn.getDeploymentID(evalID, time.Millisecond*600)
	assert.Error(t, err)
	assert.Empty(t, deplID)
}
