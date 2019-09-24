package nomadWorker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnector(t *testing.T) {

	// AWSProfile and AWSRegion not specified
	cfg := Config{}
	connector, err := cfg.New()

	assert.Nil(t, connector)
	assert.Error(t, err)

	// NomadSrvAddr missing
	cfg = Config{AWSProfile: "test"}
	connector, err = cfg.New()
	assert.Nil(t, connector)
	assert.Error(t, err)

	// Success
	cfg = Config{AWSProfile: "test", NomadServerAddress: "http://nomad.io"}
	connector, err = cfg.New()
	assert.NotNil(t, connector)
	assert.NoError(t, err)
}
