package main

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/api"
	"github.com/thomasobenaus/sokar/config"
	"github.com/thomasobenaus/sokar/helper"
	mock_logging "github.com/thomasobenaus/sokar/test/mocks/logging"
	mock_scaler "github.com/thomasobenaus/sokar/test/mocks/scaler"
)

func Test_SetupScaleScheduleShouldSucceed(t *testing.T) {
	logger := zerolog.Logger{}
	_, err := setupSchedule(nil, logger)
	assert.Error(t, err)

	// GIVEN
	entries := make([]config.ScaleScheduleEntry, 0)

	// normal entry
	start := helper.SimpleTime{Hour: 9, Minute: 30}
	end := helper.SimpleTime{Hour: 10, Minute: 30}
	entry := config.ScaleScheduleEntry{Days: []time.Weekday{time.Monday, time.Wednesday}, StartTime: start, EndTime: end, MinScale: 1, MaxScale: 30}
	entries = append(entries, entry)

	// entry with min scale unbound
	start = helper.SimpleTime{Hour: 11, Minute: 30}
	end = helper.SimpleTime{Hour: 12, Minute: 30}
	entry = config.ScaleScheduleEntry{Days: []time.Weekday{time.Monday, time.Wednesday}, StartTime: start, EndTime: end, MinScale: -1, MaxScale: 30}
	entries = append(entries, entry)

	// entry with max scale unbound
	start = helper.SimpleTime{Hour: 13, Minute: 30}
	end = helper.SimpleTime{Hour: 14, Minute: 30}
	entry = config.ScaleScheduleEntry{Days: []time.Weekday{time.Monday, time.Wednesday}, StartTime: start, EndTime: end, MinScale: 1, MaxScale: -1}
	entries = append(entries, entry)

	// conflicting entry
	start = helper.SimpleTime{Hour: 9, Minute: 35}
	end = helper.SimpleTime{Hour: 14, Minute: 30}
	entry = config.ScaleScheduleEntry{Days: []time.Weekday{time.Monday, time.Wednesday}, StartTime: start, EndTime: end, MinScale: 1, MaxScale: 10}
	entries = append(entries, entry)

	cap := config.CapacityPlanner{ScaleSchedule: entries}
	cfg := config.Config{CapacityPlanner: cap}

	// WHEN
	schedule, err := setupSchedule(&cfg, logger)

	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, schedule)
	at := helper.SimpleTime{Hour: 9, Minute: 31}
	minScale, maxScale, err := schedule.ScaleRangeAt(time.Monday, at)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), minScale)
	assert.Equal(t, uint(30), maxScale)
	minScale, maxScale, err = schedule.ScaleRangeAt(time.Wednesday, at)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), minScale)
	assert.Equal(t, uint(30), maxScale)
	_, _, err = schedule.ScaleRangeAt(time.Thursday, at)
	assert.Error(t, err)
}

func Test_CliAndConfig(t *testing.T) {

	nomadSrvAddr := "http://nomad.example.com:4646"
	cfgFile := "examples/config/full.yaml"

	args := []string{"./sokar-bin", "--config-file=" + cfgFile, "--sca.nomad.server-address=" + nomadSrvAddr}
	cfg, err := cliAndConfig(args)

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, nomadSrvAddr, cfg.Scaler.Nomad.ServerAddr)
	assert.Len(t, cfg.CapacityPlanner.ScaleSchedule, 2)

	args = []string{"./sokar-bin", "--co"}
	cfg, err = cliAndConfig(args)
	assert.Error(t, err)
	assert.Nil(t, cfg)

	args = []string{"./sokar-bin", "--sca.nomad.server-address=" + nomadSrvAddr}
	cfg, err = cliAndConfig(args)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func Test_SetupLogging(t *testing.T) {
	_, err := setupLogging(nil)
	assert.Error(t, err)

	cfg := config.Config{}
	lf, err := setupLogging(&cfg)
	assert.NoError(t, err)
	assert.NotNil(t, lf)
}

func Test_SetupScaler_Failures(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logF := mock_logging.NewMockLoggerFactory(mockCtrl)

	// no logging factory
	scaler, err := setupScaler("any", 0, 1, time.Second*1, nil, nil, false)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	scaler, err = setupScaler("any", 0, 1, time.Second*1, nil, logF, false)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	// invalid watcher-interval
	scaler, err = setupScaler("any", 0, 1, time.Second*0, nil, nil, false)
	assert.Error(t, err)
	assert.Nil(t, scaler)
}

func Test_SetupScaler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logF := mock_logging.NewMockLoggerFactory(mockCtrl)
	scalingTarget := mock_scaler.NewMockScalingTarget(mockCtrl)
	logF.EXPECT().NewNamedLogger(gomock.Any()).Times(1)

	scaler, err := setupScaler("any", 0, 1, time.Second*1, scalingTarget, logF, false)
	assert.NoError(t, err)
	assert.NotNil(t, scaler)
}

func Test_SetupScaleEmitters(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logF := mock_logging.NewMockLoggerFactory(mockCtrl)

	emitters, err := setupScaleAlertEmitters(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, emitters)

	apiInst := api.New(12000)
	emitters, err = setupScaleAlertEmitters(apiInst, nil)
	assert.Error(t, err)
	assert.Nil(t, emitters)

	logF.EXPECT().NewNamedLogger(gomock.Any())
	emitters, err = setupScaleAlertEmitters(apiInst, logF)
	assert.NoError(t, err)
	assert.Len(t, emitters, 1)
}

func Test_SetupScalingTarget(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := config.Scaler{
		Nomad: config.SCANomad{ServerAddr: "http://nomad"},
	}
	scalingTarget, err := setupScalingTarget(cfg, nil)
	assert.Error(t, err)
	assert.Nil(t, scalingTarget)

	logF := mock_logging.NewMockLoggerFactory(mockCtrl)
	logF.EXPECT().NewNamedLogger(gomock.Any()).Times(1)

	cfg = config.Scaler{}
	scalingTarget, err = setupScalingTarget(cfg, logF)
	assert.Error(t, err)
	assert.Nil(t, scalingTarget)

	cfg = config.Scaler{
		Mode: config.ScalerModeNomadDataCenter,
		Nomad: config.SCANomad{
			ServerAddr:    "http://nomad",
			DataCenterAWS: config.SCANomadDataCenterAWS{Region: "eu-central-1"},
		},
	}

	logF.EXPECT().NewNamedLogger(gomock.Any()).Times(1)
	scalingTarget, err = setupScalingTarget(cfg, logF)
	assert.NoError(t, err)
	assert.NotNil(t, scalingTarget)

	cfg = config.Scaler{
		Mode:  config.ScalerModeNomadJob,
		Nomad: config.SCANomad{ServerAddr: "http://nomad"},
	}
	logF.EXPECT().NewNamedLogger(gomock.Any()).Times(1)
	scalingTarget, err = setupScalingTarget(cfg, logF)
	assert.NoError(t, err)
	assert.NotNil(t, scalingTarget)

	cfg = config.Scaler{
		Mode:   config.ScalerModeAwsEc2,
		AwsEc2: config.SCAAwsEc2{Region: "eu-central-1", Profile: "test-profile", ASGTagKey: "key"},
	}
	logF.EXPECT().NewNamedLogger(gomock.Any()).Times(1)
	scalingTarget, err = setupScalingTarget(cfg, logF)
	assert.NoError(t, err)
	assert.NotNil(t, scalingTarget)
}
