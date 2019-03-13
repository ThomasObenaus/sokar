package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/config"
)

func Test_CliAndConfig(t *testing.T) {

	nomadSrvAddr := "http://nomad.example.com:4646"
	cfgFile := "cfg.yaml"

	args := []string{"./sokar-bin", "--config-file=" + cfgFile, "--nomad-server-address=" + nomadSrvAddr}
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
	assert.Error(t, err)
	assert.Nil(t, cfg)

	args = []string{"./sokar-bin", "--config-file=" + cfgFile}
	cfg, err = cliAndConfig(args)
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func Test_SetupLogging(t *testing.T) {
	_, err := setupLogging(nil)
	assert.Error(t, err)

	cfg := config.Config{}
	lf, err := setupLogging(&cfg)
	assert.NoError(t, err)
	assert.NotNil(t, lf)
}
