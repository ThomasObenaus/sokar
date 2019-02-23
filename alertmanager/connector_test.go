package alertmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnector(t *testing.T) {

	cfg := Config{}
	connector := cfg.New()

	assert.NotNil(t, connector)
}
