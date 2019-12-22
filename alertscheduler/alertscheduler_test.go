package alertscheduler

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewShouldCreateInstance(t *testing.T) {
	alertscheduler := New()
	assert.NotNil(t, alertscheduler)
}

func Test_WithLogger(t *testing.T) {

	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	am := New(WithLogger(logger))
	require.NotNil(t, am)
	assert.Equal(t, zerolog.DebugLevel, logger.GetLevel())
}
