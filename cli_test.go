package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseArgs(t *testing.T) {

	nomadSrvAddr := "http://nomad.example.com:4646"
	cfgFile := "cfg.yaml"

	args := []string{"./sokar-bin", "--config-file=" + cfgFile, "--nomad-server-address=" + nomadSrvAddr}
	cliArgs, err := parseArgs(args)
	require.NoError(t, err)
	assert.Equal(t, cfgFile, cliArgs.CfgFile)
	assert.Equal(t, nomadSrvAddr, cliArgs.NomadServerAddr)

	args = []string{"./sokar-bin", "--config-file=" + cfgFile}
	cliArgs, err = parseArgs(args)
	require.NoError(t, err)
	assert.Equal(t, cfgFile, cliArgs.CfgFile)
	assert.Empty(t, cliArgs.NomadServerAddr)

	args = []string{"./sokar-bin", "--unknown"}
	cliArgs, err = parseArgs(args)
	assert.Error(t, err)

	cliArgs, err = parseArgs(nil)
	assert.Error(t, err)

	cliArgs, err = parseArgs([]string{})
	assert.Error(t, err)
}

func Test_ValidateArgs(t *testing.T) {

	cfgFile := "cfg.yaml"

	args := []string{"./sokar-bin", "--config-file=" + cfgFile}
	cliArgs, err := parseArgs(args)
	require.NoError(t, err)
	assert.True(t, cliArgs.validateArgs())

	args = []string{"./sokar-bin"}
	cliArgs, err = parseArgs(args)
	require.NoError(t, err)
	assert.False(t, cliArgs.validateArgs())
}
