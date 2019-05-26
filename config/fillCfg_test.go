package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FillCfg_Flags(t *testing.T) {

	cfg := NewDefaultConfig()
	args := []string{
		"--dry-run",
		"--dummy-scaling-target",
		"--sca.mode=dc",
		"--port=1000",
		"--nomad.server-address=http://nomad",
		"--job.name=job",
		"--job.min=10",
		"--job.max=100",
		"--logging.structured",
		"--logging.unix-ts",
		"--cap.down-scale-cooldown=90s",
		"--cap.up-scale-cooldown=91s",
		"--saa.no-alert-damping=100",
		"--saa.up-thresh=101",
		"--saa.down-thresh=102",
		"--saa.eval-cycle=103s",
		"--saa.cleanup-cycle=104s",
		"--saa.eval-period-factor=105",
		"--saa.scale-alerts=alert 1:1.2:This is an upscaling alert",
		"--saa.alert-expiration-time=5m",
	}

	err := cfg.ReadConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeDataCenter, cfg.Scaler.Mode)
	assert.True(t, cfg.DummyScalingTarget)
	assert.True(t, cfg.DryRunMode)
	assert.Equal(t, 1000, cfg.Port)
	assert.Equal(t, "http://nomad", cfg.Nomad.ServerAddr)
	assert.Equal(t, "job", cfg.Job.Name)
	assert.Equal(t, uint(10), cfg.Job.MinCount)
	assert.Equal(t, uint(100), cfg.Job.MaxCount)
	assert.Equal(t, time.Duration(time.Second*90), cfg.CapacityPlanner.DownScaleCooldownPeriod)
	assert.Equal(t, time.Duration(time.Second*91), cfg.CapacityPlanner.UpScaleCooldownPeriod)
	assert.True(t, cfg.Logging.Structured)
	assert.True(t, cfg.Logging.UxTimestamp)
	assert.Equal(t, float32(100), cfg.ScaleAlertAggregator.NoAlertScaleDamping)
	assert.Equal(t, float32(101), cfg.ScaleAlertAggregator.UpScaleThreshold)
	assert.Equal(t, float32(102), cfg.ScaleAlertAggregator.DownScaleThreshold)
	assert.Equal(t, time.Duration(time.Second*103), cfg.ScaleAlertAggregator.EvaluationCycle)
	assert.Equal(t, time.Duration(time.Second*104), cfg.ScaleAlertAggregator.CleanupCycle)
	assert.Equal(t, uint(105), cfg.ScaleAlertAggregator.EvaluationPeriodFactor)

	require.Len(t, cfg.ScaleAlertAggregator.ScaleAlerts, 1)
	assert.Equal(t, "alert 1", cfg.ScaleAlertAggregator.ScaleAlerts[0].Name)
	assert.Equal(t, float32(1.2), cfg.ScaleAlertAggregator.ScaleAlerts[0].Weight)
	assert.Equal(t, "This is an upscaling alert", cfg.ScaleAlertAggregator.ScaleAlerts[0].Description)
	assert.Equal(t, time.Minute*5, cfg.ScaleAlertAggregator.AlertExpirationTime)
}

func Test_AlertMapToAlerts(t *testing.T) {

	alerts, err := alertMapToAlerts(nil)
	assert.Nil(t, alerts)
	assert.Error(t, err)

	// Success
	alert1 := make(map[string]string, 0)
	alert1["name"] = "Alert 1"
	alert1["weight"] = "2.0"
	alert1["description"] = "Alert for Upscaling"
	alertList := make([]map[string]string, 0)
	alertList = append(alertList, alert1)

	alerts, err = alertMapToAlerts(alertList)
	assert.NotNil(t, alerts)
	assert.NoError(t, err)
	assert.Equal(t, "Alert 1", alerts[0].Name)
	assert.Equal(t, "Alert for Upscaling", alerts[0].Description)
	assert.Equal(t, float32(2), alerts[0].Weight)

	// Fail - invalid weight
	alert1 = make(map[string]string, 0)
	alert1["name"] = "Alert 1"
	alert1["weight"] = "invalid weight"
	alert1["description"] = "Alert for Upscaling"
	alertList = make([]map[string]string, 0)
	alertList = append(alertList, alert1)

	alerts, err = alertMapToAlerts(alertList)
	assert.Nil(t, alerts)
	assert.Error(t, err)

	// Fail - missing name
	alert1 = make(map[string]string, 0)
	alertList = make([]map[string]string, 0)
	alertList = append(alertList, alert1)

	alerts, err = alertMapToAlerts(alertList)
	assert.Nil(t, alerts)
	assert.Error(t, err)
}
func Test_AlertStrToAlerts(t *testing.T) {
	alerts, err := alertStrToAlerts("")
	assert.Empty(t, alerts)
	assert.NoError(t, err)

	// Success
	alerts, err = alertStrToAlerts("alert 1:1.2:This is an upscaling alert;alert 2:-1.2")
	assert.NotEmpty(t, alerts)
	assert.NoError(t, err)
	assert.Len(t, alerts, 2)
	assert.Equal(t, "alert 1", alerts[0].Name)
	assert.Equal(t, float32(1.2), alerts[0].Weight)
	assert.Equal(t, "This is an upscaling alert", alerts[0].Description)
	assert.Equal(t, "alert 2", alerts[1].Name)
	assert.Equal(t, float32(-1.2), alerts[1].Weight)
	assert.Equal(t, "", alerts[1].Description)

	// Success - robust
	alerts, err = alertStrToAlerts("alert 1  :  1.2 :   This is an upscaling alert   ;  alert 2 :-1.2 : ;;;")
	assert.NotEmpty(t, alerts)
	assert.NoError(t, err)
	assert.Len(t, alerts, 2)
	assert.Equal(t, "alert 1", alerts[0].Name)
	assert.Equal(t, float32(1.2), alerts[0].Weight)
	assert.Equal(t, "This is an upscaling alert", alerts[0].Description)
	assert.Equal(t, "alert 2", alerts[1].Name)
	assert.Equal(t, float32(-1.2), alerts[1].Weight)
	assert.Equal(t, "", alerts[1].Description)

	// Fail - invalid weight
	alerts, err = alertStrToAlerts("alert 1:sldkj")
	assert.Empty(t, alerts)
	assert.Error(t, err)

	// Fail - invalid
	alerts, err = alertStrToAlerts("alert 1")
	assert.Empty(t, alerts)
	assert.Error(t, err)
}

func Test_ExtractAlertsFromViper(t *testing.T) {

	vp := viper.New()
	alerts, err := extractAlertsFromViper(vp)
	assert.Empty(t, alerts)
	assert.NoError(t, err)

	// Success - commandline
	vp.Set(saaScaleAlerts.name, "alert 1:1.2:This is an upscaling alert;alert 2:-1.2")
	alerts, err = extractAlertsFromViper(vp)
	assert.NotEmpty(t, alerts)
	assert.NoError(t, err)
	assert.Len(t, alerts, 2)
	assert.Equal(t, "alert 1", alerts[0].Name)
	assert.Equal(t, float32(1.2), alerts[0].Weight)
	assert.Equal(t, "This is an upscaling alert", alerts[0].Description)
	assert.Equal(t, "alert 2", alerts[1].Name)
	assert.Equal(t, float32(-1.2), alerts[1].Weight)
	assert.Equal(t, "", alerts[1].Description)

	// Success - config
	alert1 := make(map[string]string, 0)
	alert1["name"] = "Alert 1"
	alert1["weight"] = "2.0"
	alert1["description"] = "Alert for Upscaling"
	alertList := make([]map[string]string, 0)
	alertList = append(alertList, alert1)
	vp.Set(saaScaleAlerts.name, alertList)
	alerts, err = extractAlertsFromViper(vp)
	assert.NotEmpty(t, alerts)
	assert.NoError(t, err)
	assert.Len(t, alerts, 1)
	assert.Equal(t, "Alert 1", alerts[0].Name)
	assert.Equal(t, float32(2), alerts[0].Weight)
	assert.Equal(t, "Alert for Upscaling", alerts[0].Description)

	// Fail - config - weight
	alert1 = make(map[string]string, 0)
	alert1["name"] = "Alert 1"
	alert1["weight"] = "qwewe"
	alertList = make([]map[string]string, 0)
	alertList = append(alertList, alert1)
	vp.Set(saaScaleAlerts.name, alertList)
	alerts, err = extractAlertsFromViper(vp)
	assert.Empty(t, alerts)
	assert.Error(t, err)

	// Success - config - empty
	vp.Set(saaScaleAlerts.name, "")
	alerts, err = extractAlertsFromViper(vp)
	assert.Empty(t, alerts)
	assert.NoError(t, err)
}

func Test_StrToScalerMode(t *testing.T) {
	mode, err := strToScalerMode("")
	assert.Empty(t, mode)
	assert.Error(t, err)

	mode, err = strToScalerMode("JOB")
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeJob, mode)

	mode, err = strToScalerMode("Dc")
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeDataCenter, mode)
}
