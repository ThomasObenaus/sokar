package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()
	assert.Equal(t, 11000, config.Port)
	assert.Equal(t, ScalerModeNomadJob, config.Scaler.Mode)
	assert.Equal(t, "scale-object", config.Scaler.AwsEc2.ASGTagKey)
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
	assert.Equal(t, zerolog.InfoLevel, config.Logging.Level)
	assert.True(t, config.CapacityPlanner.ConstantMode.Enable)
	assert.False(t, config.CapacityPlanner.LinearMode.Enable)
	assert.NotNil(t, config.CapacityPlanner.ScaleSchedule)
	require.Empty(t, config.CapacityPlanner.ScaleSchedule)
}

func Test_Flags(t *testing.T) {

	args := []string{
		"--dry-run",
		"--port=1000",
		"--scale-object.name=job",
		"--scale-object.min=10",
		"--scale-object.max=100",
		"--logging.structured",
		"--logging.unix-ts",
		"--logging.level=debug",
		"--cap.down-scale-cooldown=90s",
		"--cap.up-scale-cooldown=91s",
		"--saa.no-alert-damping=100",
		"--saa.up-thresh=101",
		"--saa.down-thresh=102",
		"--saa.eval-cycle=103s",
		"--saa.cleanup-cycle=104s",
		"--saa.eval-period-factor=105",
		"--saa.scale-alerts=[{'name':'alert 1','weight':1.2,'description':'This is an upscaling alert'}]",
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

	cfg, err := New(args, "SK")
	require.NoError(t, err)
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
	assert.Equal(t, zerolog.DebugLevel, cfg.Logging.Level)
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

func Test_StrToLogLevel(t *testing.T) {

	targetType := reflect.TypeOf("")

	level, err := strToLoglevel("unknown", targetType)
	assert.Error(t, err)
	assert.Equal(t, zerolog.NoLevel, level)

	level, err = strToLoglevel("debug", targetType)
	assert.NoError(t, err)
	assert.Equal(t, zerolog.DebugLevel, level)

	level, err = strToLoglevel("info", targetType)
	assert.NoError(t, err)
	assert.Equal(t, zerolog.InfoLevel, level)

	level, err = strToLoglevel("warn", targetType)
	assert.NoError(t, err)
	assert.Equal(t, zerolog.WarnLevel, level)

	level, err = strToLoglevel("error", targetType)
	assert.NoError(t, err)
	assert.Equal(t, zerolog.ErrorLevel, level)

	level, err = strToLoglevel("fatal", targetType)
	assert.NoError(t, err)
	assert.Equal(t, zerolog.FatalLevel, level)

	level, err = strToLoglevel("off", targetType)
	assert.NoError(t, err)
	assert.Equal(t, zerolog.Disabled, level)
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
