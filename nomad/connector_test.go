package nomad

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConnector(t *testing.T) {

	connector, err := New("http://1.2.3.4")
	assert.NotNil(t, connector)
	assert.NoError(t, err)
	assert.Equal(t, time.Minute, connector.deploymentTimeOut)
	assert.Equal(t, time.Second*30, connector.evaluationTimeOut)

	connector, err = New("http://1.2.3.4", WithDeploymentTimeOut(time.Second*1234), WithEvaluationTimeOut(time.Second*5678))
	assert.NotNil(t, connector)
	assert.NoError(t, err)
	assert.Equal(t, time.Second*1234, connector.deploymentTimeOut)
	assert.Equal(t, time.Second*5678, connector.evaluationTimeOut)

	connector, err = New("")
	assert.Nil(t, connector)
	assert.Error(t, err)
}

func ExampleNew() {
	conn, err := New("http://1.2.3.4")

	if err != nil {
		log.Fatalf("Unable to create connector: %s.", err.Error())
	}

	// just to avoid the not used error
	_ = conn
}
