package nomadWorker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAWSNewSessionFromProfile(t *testing.T) {
	session, err := newAWSSessionFromProfile("invalid")
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.True(t, *session.Config.CredentialsChainVerboseErrors)
}

func TestAWSNewSession(t *testing.T) {

	session, err := newAWSSession("")
	assert.Error(t, err)
	assert.Nil(t, session)

	session, err = newAWSSession("eu-central-1")
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.True(t, *session.Config.CredentialsChainVerboseErrors)
}
