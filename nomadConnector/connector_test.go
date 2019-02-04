package nomadConnector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnector(t *testing.T) {

	cfg := Config{NomadServerAddress: "http://1.2.3.4"}
	connector, err := cfg.New()

	assert.NotNil(t, connector)
	assert.NoError(t, err)

	cfg = Config{}
	connector, err = cfg.New()

	assert.Nil(t, connector)
	assert.Error(t, err)
}
