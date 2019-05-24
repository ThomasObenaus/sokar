package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()
	assert.Equal(t, 11000, config.Port)
	assert.Equal(t, false, config.DummyScalingTarget)
	assert.Equal(t, false, config.DryRunMode)
	assert.Equal(t, float32(1), config.ScaleAlertAggregator.NoAlertScaleDamping)
	assert.Equal(t, float32(10), config.ScaleAlertAggregator.UpScaleThreshold)
	assert.Equal(t, float32(-10), config.ScaleAlertAggregator.DownScaleThreshold)
	assert.Equal(t, time.Second*1, config.ScaleAlertAggregator.EvaluationCycle)
	assert.Equal(t, uint(10), config.ScaleAlertAggregator.EvaluationPeriodFactor)
	assert.Equal(t, time.Second*60, config.ScaleAlertAggregator.CleanupCycle)
	assert.Equal(t, time.Minute*10, config.ScaleAlertAggregator.AlertExpirationTime)
	assert.NotNil(t, config.ScaleAlertAggregator.ScaleAlerts)
	assert.Empty(t, config.ScaleAlertAggregator.ScaleAlerts)
}
