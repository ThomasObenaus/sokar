package awsEc2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnector(t *testing.T) {

	// ASGTagKey not specified
	connector, err := New("")
	assert.Nil(t, connector)
	assert.Error(t, err)

	// AWSProfile and AWSRegion not specified
	connector, err = New("data-center")
	assert.Nil(t, connector)
	assert.Error(t, err)

	// success (profile)
	connector, err = New("data-center", WithAwsProfile("profile"))
	assert.NotNil(t, connector)
	assert.NoError(t, err)
	assert.Equal(t, "profile", connector.awsProfile)

	// success (region)
	connector, err = New("data-center", WithAwsRegion("region"))
	assert.NotNil(t, connector)
	assert.NoError(t, err)
	assert.Equal(t, "region", connector.awsRegion)
}

func TestValidate(t *testing.T) {

	// all missing
	err := validate(Connector{})
	assert.Error(t, err)

	// region and profile missing
	err = validate(Connector{tagKey: "tagkey"})
	assert.Error(t, err)

	// success
	err = validate(Connector{tagKey: "tagkey", awsProfile: "profile"})
	assert.NoError(t, err)
}
