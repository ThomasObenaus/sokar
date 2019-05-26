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
	session, err := newAWSSession()
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.True(t, *session.Config.CredentialsChainVerboseErrors)
}
