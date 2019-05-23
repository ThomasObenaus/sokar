package nomadworker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnector(t *testing.T) {

	cfg := Config{}
	connector, err := cfg.New()

	assert.NotNil(t, connector)
	assert.NoError(t, err)
}
