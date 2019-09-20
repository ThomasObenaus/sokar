package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAWSNewSessionFromProfile(t *testing.T) {
	session, err := NewAWSSessionFromProfile("invalid")
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.True(t, *session.Config.CredentialsChainVerboseErrors)
}

func TestAWSNewSession(t *testing.T) {

	session, err := NewAWSSession("")
	assert.Error(t, err)
	assert.Nil(t, session)

	session, err = NewAWSSession("eu-central-1")
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.True(t, *session.Config.CredentialsChainVerboseErrors)
}
