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

	cfg = Config{AWSProfile: "test", ASGTagKey: "data-center"}
	connector, err = cfg.New()
	assert.NotNil(t, connector)
	assert.NoError(t, err)

	cfg = Config{AWSRegion: "test-region", ASGTagKey: "data-center"}
	connector, err = cfg.New()
	assert.NotNil(t, connector)
	assert.NoError(t, err)

	cfg = Config{AWSRegion: "test-region"}
	connector, err = cfg.New()
	assert.Nil(t, connector)
	assert.Error(t, err)
}
