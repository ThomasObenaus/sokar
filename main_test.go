package main

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/api"
	"github.com/thomasobenaus/sokar/config"
	mock_logging "github.com/thomasobenaus/sokar/test/logging"
	mock_scaler "github.com/thomasobenaus/sokar/test/scaler"
)

func Test_CliAndConfig(t *testing.T) {

	nomadSrvAddr := "http://nomad.example.com:4646"
	cfgFile := "examples/config/full.yaml"

	args := []string{"./sokar-bin", "--config-file=" + cfgFile, "--sca.nomad.server-address=" + nomadSrvAddr}
	cfg, err := cliAndConfig(args)

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, nomadSrvAddr, cfg.Scaler.Nomad.ServerAddr)

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
	scaler, err := setupScaler("any", 0, 1, time.Second*1, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	scaler, err = setupScaler("any", 0, 1, time.Second*1, nil, logF)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	// invalid watcher-interval
	scaler, err = setupScaler("any", 0, 1, time.Second*0, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, scaler)
}

func Test_SetupScaler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logF := mock_logging.NewMockLoggerFactory(mockCtrl)
	scalingTarget := mock_scaler.NewMockScalingTarget(mockCtrl)
	logF.EXPECT().NewNamedLogger(gomock.Any()).Times(1)

	scaler, err := setupScaler("any", 0, 1, time.Second*1, scalingTarget, logF)
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
}
