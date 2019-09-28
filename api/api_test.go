package api

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	api := New(12345)
	assert.NotNil(t, api)
}

func TestWithLogger(t *testing.T) {

	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	api := New(1234, WithLogger(logger))
	require.NotNil(t, api)
	assert.Equal(t, zerolog.DebugLevel, logger.GetLevel())
}

func TestRunJoinStop(t *testing.T) {

	api := New(1234)
	require.NotNil(t, api)

	api.Run()
	start := time.Now()
	api.Stop()
	api.Join()

	assert.WithinDuration(t, start.Add(time.Millisecond*500), time.Now(), time.Second*1)
}
