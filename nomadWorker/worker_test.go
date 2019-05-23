package nomadWorker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetJobCount(t *testing.T) {

	cfg := Config{}
	connector, err := cfg.New()

	require.NotNil(t, connector)
	require.NoError(t, err)

	err = connector.SetJobCount("public-services", 10)
	assert.NoError(t, err)
}

func TestGetJobCount(t *testing.T) {

	cfg := Config{}
	connector, err := cfg.New()

	require.NotNil(t, connector)
	require.NoError(t, err)

	count, err := connector.GetJobCount("public-services")
	assert.NoError(t, err)
	assert.Equal(t, uint(100), count)

	connector.SetJobCount("public-services", 10)
	count, err = connector.GetJobCount("public-services")
	assert.NoError(t, err)
	assert.Equal(t, uint(10), count)
}

func TestIsJobDead(t *testing.T) {

	cfg := Config{}
	connector, err := cfg.New()

	require.NotNil(t, connector)
	require.NoError(t, err)

	dead, err := connector.IsJobDead("public-services")
	assert.NoError(t, err)
	assert.False(t, dead)
}
