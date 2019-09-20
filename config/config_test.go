package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()
	assert.Equal(t, 11000, config.Port)
	assert.Equal(t, ScalerModeNomadJob, config.Scaler.Mode)
	assert.Equal(t, time.Duration(time.Second*5), config.Scaler.WatcherInterval)
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
	assert.Equal(t, time.Second*80, config.CapacityPlanner.DownScaleCooldownPeriod)
	assert.Equal(t, time.Second*60, config.CapacityPlanner.UpScaleCooldownPeriod)
	assert.Equal(t, uint(1), config.CapacityPlanner.ConstantMode.Offset)
	assert.True(t, config.CapacityPlanner.ConstantMode.Enable)
	assert.False(t, config.CapacityPlanner.LinearMode.Enable)
}
