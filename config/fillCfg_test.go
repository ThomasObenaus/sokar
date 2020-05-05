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
		"--sca.aws-ec2.asg-tag-key=asg-tag-key",
		"--sca.nomad.dc-aws.instance-termination-timeout=124s",
		"--sca.watcher-interval=50s",
		"--cap.constant-mode.enable=false",
		"--cap.constant-mode.offset=106",
		"--cap.linear-mode.enable=true",
		"--cap.linear-mode.scale-factor-weight=0.107",
		"--cap.scale-schedule=MON-FRI 7 9 10-30|WED-SAT 13:15 17:25 2-22",
	}

	err := cfg.ReadConfig(args)
	assert.NoError(t, err)
	assert.Equal(t, ScalerModeAwsEc2, cfg.Scaler.Mode)
	assert.Equal(t, "profile-test", cfg.Scaler.Nomad.DataCenterAWS.Profile)
	assert.Equal(t, "region-test", cfg.Scaler.Nomad.DataCenterAWS.Region)
	assert.Equal(t, time.Duration(time.Second*124), cfg.Scaler.Nomad.DataCenterAWS.InstanceTerminationTimeout)
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

	require.NotNil(t, cfg.CapacityPlanner.ScaleSchedule)
	require.Len(t, cfg.CapacityPlanner.ScaleSchedule, 2)
	assert.Len(t, cfg.CapacityPlanner.ScaleSchedule[0].Days, 5)
	assert.Equal(t, uint(7), cfg.CapacityPlanner.ScaleSchedule[0].StartTime.Hour)
	assert.Equal(t, uint(0), cfg.CapacityPlanner.ScaleSchedule[0].StartTime.Minute)
	assert.Equal(t, uint(9), cfg.CapacityPlanner.ScaleSchedule[0].EndTime.Hour)
	assert.Equal(t, uint(0), cfg.CapacityPlanner.ScaleSchedule[0].EndTime.Minute)
	assert.Equal(t, 10, cfg.CapacityPlanner.ScaleSchedule[0].MinScale)
	assert.Equal(t, 30, cfg.CapacityPlanner.ScaleSchedule[0].MaxScale)
	assert.Len(t, cfg.CapacityPlanner.ScaleSchedule[1].Days, 4)
	assert.Equal(t, uint(13), cfg.CapacityPlanner.ScaleSchedule[1].StartTime.Hour)
	assert.Equal(t, uint(15), cfg.CapacityPlanner.ScaleSchedule[1].StartTime.Minute)
	assert.Equal(t, uint(17), cfg.CapacityPlanner.ScaleSchedule[1].EndTime.Hour)
	assert.Equal(t, uint(25), cfg.CapacityPlanner.ScaleSchedule[1].EndTime.Minute)
	assert.Equal(t, 2, cfg.CapacityPlanner.ScaleSchedule[1].MinScale)
	assert.Equal(t, 22, cfg.CapacityPlanner.ScaleSchedule[1].MaxScale)
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
	alert1 := make(map[string]string)
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
	alert1 = make(map[string]string)
	alert1["name"] = "Alert 1"
	alert1["weight"] = "invalid weight"
	alert1["description"] = "Alert for Upscaling"
	alertList = make([]map[string]string, 0)
	alertList = append(alertList, alert1)

	alerts, err = alertMapToAlerts(alertList)
	assert.Nil(t, alerts)
	assert.Error(t, err)

	// Fail - missing name
	alert1 = make(map[string]string)
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
	alert1 := make(map[string]string)
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
	alert1 = make(map[string]string)
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

func Test_ValidateCapacityPlanner(t *testing.T) {
	capacityPlanner := CapacityPlanner{}
	err := validateCapacityPlanner(capacityPlanner)
	assert.Error(t, err)

	capacityPlanner = CapacityPlanner{ConstantMode: CAPConstMode{Enable: true}, LinearMode: CAPLinearMode{Enable: true}}
	err = validateCapacityPlanner(capacityPlanner)
	assert.Error(t, err)

	capacityPlanner = CapacityPlanner{LinearMode: CAPLinearMode{Enable: true, ScaleFactorWeight: 0}}
	err = validateCapacityPlanner(capacityPlanner)
	assert.Error(t, err)

	capacityPlanner = CapacityPlanner{ConstantMode: CAPConstMode{Enable: true, Offset: 0}}
	err = validateCapacityPlanner(capacityPlanner)
	assert.Error(t, err)

	capacityPlanner = CapacityPlanner{LinearMode: CAPLinearMode{Enable: true, ScaleFactorWeight: 0.5}}
	err = validateCapacityPlanner(capacityPlanner)
	assert.NoError(t, err)

	capacityPlanner = CapacityPlanner{ConstantMode: CAPConstMode{Enable: true, Offset: 1}}
	err = validateCapacityPlanner(capacityPlanner)
	assert.NoError(t, err)

	// success - no schedule
	scaleSchedule := make([]ScaleScheduleEntry, 0)
	capacityPlanner = CapacityPlanner{ConstantMode: CAPConstMode{Enable: true, Offset: 1}, ScaleSchedule: scaleSchedule}
	err = validateCapacityPlanner(capacityPlanner)
	assert.NoError(t, err)
}

func Test_ExtractScaleScheduleFromViper(t *testing.T) {

	vp := viper.New()
	scaleScheduleEntries, err := extractScaleScheduleFromViper(vp)
	assert.Empty(t, scaleScheduleEntries)
	assert.NoError(t, err)

	// Success - commandline
	vp.Set(capScaleSchedule.name, "MON-FRI 7 9 10-30|WED-SAT 13:15 17:25 2-*")
	scaleScheduleEntries, err = extractScaleScheduleFromViper(vp)
	assert.NotEmpty(t, scaleScheduleEntries)
	assert.NoError(t, err)
	assert.Len(t, scaleScheduleEntries, 2)
	assert.Len(t, scaleScheduleEntries[0].Days, 5)
	assert.Equal(t, uint(7), scaleScheduleEntries[0].StartTime.Hour)
	assert.Equal(t, uint(0), scaleScheduleEntries[0].StartTime.Minute)
	assert.Equal(t, uint(9), scaleScheduleEntries[0].EndTime.Hour)
	assert.Equal(t, uint(0), scaleScheduleEntries[0].EndTime.Minute)
	assert.Equal(t, 10, scaleScheduleEntries[0].MinScale)
	assert.Equal(t, 30, scaleScheduleEntries[0].MaxScale)
	assert.Len(t, scaleScheduleEntries[1].Days, 4)
	assert.Equal(t, uint(13), scaleScheduleEntries[1].StartTime.Hour)
	assert.Equal(t, uint(15), scaleScheduleEntries[1].StartTime.Minute)
	assert.Equal(t, uint(17), scaleScheduleEntries[1].EndTime.Hour)
	assert.Equal(t, uint(25), scaleScheduleEntries[1].EndTime.Minute)
	assert.Equal(t, 2, scaleScheduleEntries[1].MinScale)
	assert.Equal(t, -1, scaleScheduleEntries[1].MaxScale)

	// Success - config
	entries := make([]map[string]string, 0)
	entry := make(map[string]string)
	entry["days"] = "MON-FRI"
	entry["start-time"] = "7:30"
	entry["end-time"] = "9:30"
	entry["min"] = "10"
	entry["max"] = "30"
	entries = append(entries, entry)

	entry = make(map[string]string)
	entry["days"] = "SAT-SUN"
	entry["start-time"] = "17"
	entry["end-time"] = "18"
	entry["min"] = "2"
	entry["max"] = "3"
	entries = append(entries, entry)

	vp.Set(capScaleSchedule.name, entries)
	scaleScheduleEntries, err = extractScaleScheduleFromViper(vp)
	assert.NotEmpty(t, scaleScheduleEntries)
	assert.NoError(t, err)
	assert.Len(t, scaleScheduleEntries, 2)
	assert.Len(t, scaleScheduleEntries[0].Days, 5)
	assert.Equal(t, uint(7), scaleScheduleEntries[0].StartTime.Hour)
	assert.Equal(t, uint(30), scaleScheduleEntries[0].StartTime.Minute)
	assert.Equal(t, uint(9), scaleScheduleEntries[0].EndTime.Hour)
	assert.Equal(t, uint(30), scaleScheduleEntries[0].EndTime.Minute)
	assert.Equal(t, 10, scaleScheduleEntries[0].MinScale)
	assert.Equal(t, 30, scaleScheduleEntries[0].MaxScale)
	assert.Len(t, scaleScheduleEntries[1].Days, 2)
	assert.Equal(t, uint(17), scaleScheduleEntries[1].StartTime.Hour)
	assert.Equal(t, uint(0), scaleScheduleEntries[1].StartTime.Minute)
	assert.Equal(t, uint(18), scaleScheduleEntries[1].EndTime.Hour)
	assert.Equal(t, uint(0), scaleScheduleEntries[1].EndTime.Minute)
	assert.Equal(t, 2, scaleScheduleEntries[1].MinScale)
	assert.Equal(t, 3, scaleScheduleEntries[1].MaxScale)

	// Fail - config - schedule
	entries = make([]map[string]string, 0)
	entry = make(map[string]string)
	entry["days"] = "invalid"
	entries = append(entries, entry)
	vp.Set(capScaleSchedule.name, entries)
	scaleScheduleEntries, err = extractScaleScheduleFromViper(vp)
	assert.Empty(t, scaleScheduleEntries)
	assert.Error(t, err)
	//
	// Success - config - empty
	vp.Set(capScaleSchedule.name, "")
	scaleScheduleEntries, err = extractScaleScheduleFromViper(vp)
	assert.Empty(t, scaleScheduleEntries)
	assert.NoError(t, err)
}

func Test_ScaleScheduleMapToScaleSchedule(t *testing.T) {
	// TODO: Reenable test
	scaleScheduleEntries, err := scaleScheduleMapToScaleSchedule(nil)
	assert.Nil(t, scaleScheduleEntries)
	assert.Error(t, err)

	// Success
	entries := make([]map[string]string, 0)
	entry := make(map[string]string)
	entry["days"] = "MON-FRI"
	entry["start-time"] = "7:30"
	entry["end-time"] = "9:30"
	entry["min"] = "10"
	entry["max"] = "30"
	entries = append(entries, entry)
	scaleScheduleEntries, err = scaleScheduleMapToScaleSchedule(entries)
	assert.NoError(t, err)
	assert.Len(t, scaleScheduleEntries[0].Days, 5)
	assert.Equal(t, uint(7), scaleScheduleEntries[0].StartTime.Hour)
	assert.Equal(t, uint(30), scaleScheduleEntries[0].StartTime.Minute)
	assert.Equal(t, uint(9), scaleScheduleEntries[0].EndTime.Hour)
	assert.Equal(t, uint(30), scaleScheduleEntries[0].EndTime.Minute)
	assert.Equal(t, 10, scaleScheduleEntries[0].MinScale)
	assert.Equal(t, 30, scaleScheduleEntries[0].MaxScale)

	// Fail - config - days
	entries = make([]map[string]string, 0)
	entry = make(map[string]string)
	entry["days"] = "invalid"
	entries = append(entries, entry)
	scaleScheduleEntries, err = scaleScheduleMapToScaleSchedule(entries)
	assert.Empty(t, scaleScheduleEntries)
	assert.Error(t, err)

	// Fail - config - min
	entries = make([]map[string]string, 0)
	entry = make(map[string]string)
	entry["days"] = "MON-FRI"
	entry["start-time"] = "7:30"
	entry["end-time"] = "9:30"
	entry["min"] = "invalid"
	entry["max"] = "30"
	entries = append(entries, entry)
	scaleScheduleEntries, err = scaleScheduleMapToScaleSchedule(entries)
	assert.Empty(t, scaleScheduleEntries)
	assert.Error(t, err)

	// Fail - config - max
	entries = make([]map[string]string, 0)
	entry = make(map[string]string)
	entry["days"] = "MON-FRI"
	entry["start-time"] = "7:30"
	entry["end-time"] = "9:30"
	entry["min"] = "10"
	entry["max"] = "invalid"
	entries = append(entries, entry)
	scaleScheduleEntries, err = scaleScheduleMapToScaleSchedule(entries)
	assert.Empty(t, scaleScheduleEntries)
	assert.Error(t, err)

	// Fail - config - start time empty
	entries = make([]map[string]string, 0)
	entry = make(map[string]string)
	entry["days"] = "MON-FRI"
	entry["start-time"] = ""
	entry["end-time"] = "9:30"
	entry["min"] = "10"
	entry["max"] = "30"
	entries = append(entries, entry)
	scaleScheduleEntries, err = scaleScheduleMapToScaleSchedule(entries)
	assert.Empty(t, scaleScheduleEntries)
	assert.Error(t, err)

	// Fail - config - end time empty
	entries = make([]map[string]string, 0)
	entry = make(map[string]string)
	entry["days"] = "MON-FRI"
	entry["start-time"] = "7:30"
	entry["end-time"] = ""
	entry["min"] = "10"
	entry["max"] = "30"
	entries = append(entries, entry)
	scaleScheduleEntries, err = scaleScheduleMapToScaleSchedule(entries)
	assert.Empty(t, scaleScheduleEntries)
	assert.Error(t, err)
}
