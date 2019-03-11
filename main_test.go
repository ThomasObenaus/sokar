package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CliAndConfig(t *testing.T) {

	nomadSrvAddr := "http://nomad.example.com:4646"
	cfgFile := "cfg.yaml"

	args := []string{"./sokar-bin", "--config-file=" + cfgFile, "--nomad-server-address=" + nomadSrvAddr}
	cfg, err := cliAndConfig(args)

	require.NoError(t, err)
	assert.Equal(t, nomadSrvAddr, cfg.Nomad.ServerAddr)

	args = []string{"./sokar-bin", "--co"}
	cfg, err = cliAndConfig(args)
	assert.Error(t, err)

	args = []string{"./sokar-bin"}
	cfg, err = cliAndConfig(args)
	assert.Error(t, err)

	args = []string{"./sokar-bin", "--config-file=" + cfgFile}
	cfg, err = cliAndConfig(args)
	assert.Error(t, err)
}
