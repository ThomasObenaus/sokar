package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/api"
	"github.com/thomasobenaus/sokar/config"
	"github.com/thomasobenaus/sokar/test/logging"
)

func Test_CliAndConfig(t *testing.T) {

	nomadSrvAddr := "http://nomad.example.com:4646"
	cfgFile := "examples/config/full.yaml"

	args := []string{"./sokar-bin", "--config-file=" + cfgFile, "--nomad.server-address=" + nomadSrvAddr}
	cfg, err := cliAndConfig(args)

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, nomadSrvAddr, cfg.Nomad.ServerAddr)

	args = []string{"./sokar-bin", "--co"}
	cfg, err = cliAndConfig(args)
	assert.Error(t, err)
	assert.Nil(t, cfg)

	args = []string{"./sokar-bin"}
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

func Test_SetupScaler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logF := mock_logging.NewMockLoggerFactory(mockCtrl)

	// no logging factory
	scaler, err := setupScaler("any", 0, 1, "nomad-addr", nil)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	logF.EXPECT().NewNamedLogger(gomock.Any()).Times(1)
	scaler, err = setupScaler("any", 0, 1, "", logF)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	logF.EXPECT().NewNamedLogger(gomock.Any()).Times(2)
	scaler, err = setupScaler("any", 0, 1, "https://nomad.com", logF)
	assert.NoError(t, err)
	assert.NotNil(t, scaler)
}
func Test_SetupScaleEmitters(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logF := mock_logging.NewMockLoggerFactory(mockCtrl)
	logger := zerolog.Logger{}

	emitters, err := setupScaleAlertEmitters(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, emitters)

	apiInst := api.New(12000, logger)
	emitters, err = setupScaleAlertEmitters(apiInst, nil)
	assert.Error(t, err)
	assert.Nil(t, emitters)

	logF.EXPECT().NewNamedLogger(gomock.Any())
	emitters, err = setupScaleAlertEmitters(apiInst, logF)
	assert.NoError(t, err)
	assert.Len(t, emitters, 1)
}
