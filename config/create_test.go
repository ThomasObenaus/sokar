package config

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davecgh/go-spew/spew"
)

var fullConfig = `
port: 1234
dry_run_mode: true
nomad:
  srv_addr: "http://localhost:4646"
job:
  name: "fail-service"
  min: 1
  max: 10
capacity_planner:
  down_scale_cooldown: 20s
  up_scale_cooldown: 10s
scale_alert_aggregator:
  no_alert_damping: 1.0
  up_thresh: 10.0
  down_thresh: -10.0
  eval_cycle: 1s
  eval_period_factor: 10
  cleanup_cycle: 60s
  scale_alerts:
    - name: "AlertA"
      weight: 1.5
    - name: "AlertB"
      weight: -1.5
      description: "Down alert"
logging:
  structured: false
  unix_ts: false
`
var minimalConfig = `
nomad:
  srv_addr: "http://localhost:4646"
job:
  name: "fail-service"
  min: 1
  max: 10
scale_alert_aggregator:
  scale_alerts:
    - name: "AlertA"
      weight: 1.5
    - name: "AlertB"
      weight: -1.5
      description: "Down alert"
`

var invalidConfig = `
invalid:
:yaml
`

func Test_NewconfigFromYAMLFile(t *testing.T) {
	config, err := NewConfigFromYAMLFile("NO_FILE")
	assert.Error(t, err)

	config, err = NewConfigFromYAMLFile("../test/config/full.yaml")
	assert.NoError(t, err)

	// dry_run_mode
	assert.Equal(t, false, config.DryRunMode)

	// port
	assert.Equal(t, 11000, config.Port)

	// nomad
	assert.Equal(t, "http://localhost:4646", config.Nomad.ServerAddr)

	// logging
	assert.False(t, config.Logging.Structured)
	assert.False(t, config.Logging.UxTimestamp)

	// job
	assert.Equal(t, "fail-service", config.Job.Name)
	assert.Equal(t, uint(1), config.Job.MinCount)
	assert.Equal(t, uint(10), config.Job.MaxCount)

	// cfg
	assert.Equal(t, float32(1), config.ScaleAlertAggregator.NoAlertScaleDamping)
	assert.Equal(t, float32(10), config.ScaleAlertAggregator.UpScaleThreshold)
	assert.Equal(t, float32(-10), config.ScaleAlertAggregator.DownScaleThreshold)
	assert.Equal(t, time.Second*1, config.ScaleAlertAggregator.EvaluationCycle)
	assert.Equal(t, uint(10), config.ScaleAlertAggregator.EvaluationPeriodFactor)
	assert.Equal(t, time.Second*60, config.ScaleAlertAggregator.CleanupCycle)

	// scale_alerts
	assert.Len(t, config.ScaleAlertAggregator.ScaleAlerts, 2)
	assert.Equal(t, "AlertA", config.ScaleAlertAggregator.ScaleAlerts[0].Name)
	assert.Equal(t, float32(1.5), config.ScaleAlertAggregator.ScaleAlerts[0].Weight)
	assert.Equal(t, "", config.ScaleAlertAggregator.ScaleAlerts[0].Description)
	assert.Equal(t, "AlertB", config.ScaleAlertAggregator.ScaleAlerts[1].Name)
	assert.Equal(t, float32(-1.5), config.ScaleAlertAggregator.ScaleAlerts[1].Weight)
	assert.Equal(t, "Down alert", config.ScaleAlertAggregator.ScaleAlerts[1].Description)

	// capacity_planner
	assert.Equal(t, time.Second*20, config.CapacityPlanner.DownScaleCooldownPeriod)
	assert.Equal(t, time.Second*10, config.CapacityPlanner.UpScaleCooldownPeriod)
}
func Test_NewconfigFromYAML_Invalid(t *testing.T) {
	reader := strings.NewReader(invalidConfig)

	_, err := NewConfigFromYAML(reader)
	assert.Error(t, err)
}

func Test_NewConfigFromYAML_Partial(t *testing.T) {
	reader := strings.NewReader(minimalConfig)

	config, err := NewConfigFromYAML(reader)
	require.NoError(t, err)

	// dry_run_mode
	assert.Equal(t, false, config.DryRunMode)

	// port
	assert.Equal(t, 11000, config.Port)

	// nomad
	assert.Equal(t, "http://localhost:4646", config.Nomad.ServerAddr)

	// logging
	assert.False(t, config.Logging.Structured)
	assert.False(t, config.Logging.UxTimestamp)

	// job
	assert.Equal(t, "fail-service", config.Job.Name)
	assert.Equal(t, uint(1), config.Job.MinCount)
	assert.Equal(t, uint(10), config.Job.MaxCount)

	// cfg
	assert.Equal(t, float32(1), config.ScaleAlertAggregator.NoAlertScaleDamping)
	assert.Equal(t, float32(10), config.ScaleAlertAggregator.UpScaleThreshold)
	assert.Equal(t, float32(-10), config.ScaleAlertAggregator.DownScaleThreshold)
	assert.Equal(t, time.Second*1, config.ScaleAlertAggregator.EvaluationCycle)
	assert.Equal(t, uint(10), config.ScaleAlertAggregator.EvaluationPeriodFactor)
	assert.Equal(t, time.Second*60, config.ScaleAlertAggregator.CleanupCycle)

	// scale_alerts
	assert.Len(t, config.ScaleAlertAggregator.ScaleAlerts, 2)
	assert.Equal(t, "AlertA", config.ScaleAlertAggregator.ScaleAlerts[0].Name)
	assert.Equal(t, float32(1.5), config.ScaleAlertAggregator.ScaleAlerts[0].Weight)
	assert.Equal(t, "", config.ScaleAlertAggregator.ScaleAlerts[0].Description)
	assert.Equal(t, "AlertB", config.ScaleAlertAggregator.ScaleAlerts[1].Name)
	assert.Equal(t, float32(-1.5), config.ScaleAlertAggregator.ScaleAlerts[1].Weight)
	assert.Equal(t, "Down alert", config.ScaleAlertAggregator.ScaleAlerts[1].Description)

	// capacity_planner
	assert.Equal(t, time.Second*80, config.CapacityPlanner.DownScaleCooldownPeriod)
	assert.Equal(t, time.Second*60, config.CapacityPlanner.UpScaleCooldownPeriod)
}

func Test_NewConfigFromYAML_Full(t *testing.T) {

	reader := strings.NewReader(fullConfig)

	config, err := NewConfigFromYAML(reader)
	require.NoError(t, err)

	// nomad
	assert.Equal(t, "http://localhost:4646", config.Nomad.ServerAddr)

	// logging
	assert.False(t, config.Logging.Structured)
	assert.False(t, config.Logging.UxTimestamp)

	// job
	assert.Equal(t, "fail-service", config.Job.Name)
	assert.Equal(t, uint(1), config.Job.MinCount)
	assert.Equal(t, uint(10), config.Job.MaxCount)

	// cfg
	assert.Equal(t, float32(1), config.ScaleAlertAggregator.NoAlertScaleDamping)
	assert.Equal(t, float32(10), config.ScaleAlertAggregator.UpScaleThreshold)
	assert.Equal(t, float32(-10), config.ScaleAlertAggregator.DownScaleThreshold)
	assert.Equal(t, time.Second*1, config.ScaleAlertAggregator.EvaluationCycle)
	assert.Equal(t, uint(10), config.ScaleAlertAggregator.EvaluationPeriodFactor)
	assert.Equal(t, time.Second*60, config.ScaleAlertAggregator.CleanupCycle)

	// scale_alerts
	assert.Len(t, config.ScaleAlertAggregator.ScaleAlerts, 2)
	assert.Equal(t, "AlertA", config.ScaleAlertAggregator.ScaleAlerts[0].Name)
	assert.Equal(t, float32(1.5), config.ScaleAlertAggregator.ScaleAlerts[0].Weight)
	assert.Equal(t, "", config.ScaleAlertAggregator.ScaleAlerts[0].Description)
	assert.Equal(t, "AlertB", config.ScaleAlertAggregator.ScaleAlerts[1].Name)
	assert.Equal(t, float32(-1.5), config.ScaleAlertAggregator.ScaleAlerts[1].Weight)
	assert.Equal(t, "Down alert", config.ScaleAlertAggregator.ScaleAlerts[1].Description)

	// capacity_planner
	assert.Equal(t, time.Second*20, config.CapacityPlanner.DownScaleCooldownPeriod)
	assert.Equal(t, time.Second*10, config.CapacityPlanner.UpScaleCooldownPeriod)

	spew.Dump(config)
}

func Test_NewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()
	assert.Equal(t, 11000, config.Port)
	assert.Equal(t, false, config.DryRunMode)
	assert.Equal(t, float32(1), config.ScaleAlertAggregator.NoAlertScaleDamping)
	assert.Equal(t, float32(10), config.ScaleAlertAggregator.UpScaleThreshold)
	assert.Equal(t, float32(-10), config.ScaleAlertAggregator.DownScaleThreshold)
	assert.Equal(t, time.Second*1, config.ScaleAlertAggregator.EvaluationCycle)
	assert.Equal(t, uint(10), config.ScaleAlertAggregator.EvaluationPeriodFactor)
	assert.Equal(t, time.Second*60, config.ScaleAlertAggregator.CleanupCycle)
	assert.NotNil(t, config.ScaleAlertAggregator.ScaleAlerts)
	assert.Empty(t, config.ScaleAlertAggregator.ScaleAlerts)
}
