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
		"--port=1000",
		"--scale-object.name=job",
		"--scale-object.min=10",
		"--scale-object.max=100",
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
		"--sca.mode=aws-ec2",
		"--sca.nomad.server-address=http://nomad",
		"--sca.nomad.dc-aws.region=region-test",
		"--sca.nomad.dc-aws.profile=profile-test",
		"--sca.aws-ec2.profile=profile-test",
		"--sca.aws-ec2.region=region-test",
		"--sca.aws-ec2.asg-tag-key=asg-tag-key",
		"--sca.watcher-interval=50s",
		"--cap.constant-mode.enable=false",
		"--cap.constant-mode.offset=106",
		"--cap.linear-mode.enable=true",
		"--cap.linear-mode.scale-factor-weight=0.107",
	}

	err := cfg.ReadConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeAwsEc2, cfg.Scaler.Mode)
	assert.Equal(t, "profile-test", cfg.Scaler.Nomad.DataCenterAWS.Profile)
	assert.Equal(t, "region-test", cfg.Scaler.Nomad.DataCenterAWS.Region)
	assert.Equal(t, "http://nomad", cfg.Scaler.Nomad.ServerAddr)
	assert.Equal(t, "profile-test", cfg.Scaler.AwsEc2.Profile)
	assert.Equal(t, "region-test", cfg.Scaler.AwsEc2.Region)
	assert.Equal(t, "asg-tag-key", cfg.Scaler.AwsEc2.ASGTagKey)
	assert.Equal(t, time.Duration(time.Second*50), cfg.Scaler.WatcherInterval)
	assert.True(t, cfg.DryRunMode)
	assert.Equal(t, 1000, cfg.Port)
	assert.Equal(t, "job", cfg.ScaleObject.Name)
	assert.Equal(t, uint(10), cfg.ScaleObject.MinCount)
	assert.Equal(t, uint(100), cfg.ScaleObject.MaxCount)
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
	assert.False(t, cfg.CapacityPlanner.ConstantMode.Enable)
	assert.Equal(t, uint(106), cfg.CapacityPlanner.ConstantMode.Offset)
	assert.True(t, cfg.CapacityPlanner.LinearMode.Enable)
	assert.Equal(t, float64(0.107), cfg.CapacityPlanner.LinearMode.ScaleFactorWeight)
}

func Test_ValidateScaler_NomadJob(t *testing.T) {
	sca := Scaler{}
	err := validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeNomadJob}
	err = validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeNomadJob, Nomad: SCANomad{ServerAddr: "http://test.com"}, WatcherInterval: time.Millisecond * 499}
	err = validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeNomadJob, Nomad: SCANomad{ServerAddr: "http://test.com"}, WatcherInterval: time.Second * 5}
	err = validateScaler(sca)
	assert.NoError(t, err)
}

func Test_ValidateScaler_AwsEc2(t *testing.T) {
	sca := Scaler{}
	err := validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeAwsEc2}
	err = validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeAwsEc2, AwsEc2: SCAAwsEc2{Profile: "profile"}}
	err = validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeAwsEc2, AwsEc2: SCAAwsEc2{Profile: "profile", Region: "test-region"}}
	err = validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeAwsEc2, AwsEc2: SCAAwsEc2{Profile: "profile", ASGTagKey: "datacenter"}}
	err = validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeAwsEc2, AwsEc2: SCAAwsEc2{Profile: "profile", ASGTagKey: "datacenter"}, WatcherInterval: time.Second * 2}
	err = validateScaler(sca)
	assert.NoError(t, err)
}

func Test_ValidateScaler_NomadDatacenter(t *testing.T) {
	sca := Scaler{}
	err := validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeNomadDataCenter}
	err = validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeNomadDataCenter, Nomad: SCANomad{ServerAddr: "http://test.com"}}
	err = validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeNomadDataCenter, Nomad: SCANomad{ServerAddr: "http://test.com", DataCenterAWS: SCANomadDataCenterAWS{Profile: "profile"}}}
	err = validateScaler(sca)
	assert.Error(t, err)

	sca = Scaler{Mode: ScalerModeNomadDataCenter, WatcherInterval: time.Second * 2, Nomad: SCANomad{ServerAddr: "http://test.com", DataCenterAWS: SCANomadDataCenterAWS{Profile: "profile"}}}
	err = validateScaler(sca)
	assert.NoError(t, err)
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
	assert.Equal(t, ScalerModeNomadJob, mode)

	mode, err = strToScalerMode("Dc")
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeNomadDataCenter, mode)

	mode, err = strToScalerMode("aws-eC2")
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeAwsEc2, mode)
}

// TODO: Remove as soon as the sca.nomad.mode flag has been removed
func Test_FillCfg_SupportDeprecatedFlags(t *testing.T) {

	cfg := NewDefaultConfig()
	args := []string{
		"--sca.nomad.mode=aws-ec2",
		"--sca.mode=job",
	}

	err := cfg.ReadConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeAwsEc2, cfg.Scaler.Mode)

	// old job
	cfg = NewDefaultConfig()
	args = []string{
		"--sca.mode=job",
	}

	err = cfg.ReadConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeNomadJob, cfg.Scaler.Mode)

	// old dc
	cfg = NewDefaultConfig()
	args = []string{
		"--sca.mode=dc",
	}

	err = cfg.ReadConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeNomadDataCenter, cfg.Scaler.Mode)
}
