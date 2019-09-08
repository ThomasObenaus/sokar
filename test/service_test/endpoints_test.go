package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/config"
)

func Test_EndPoints_ConfigEndpoint(t *testing.T) {
	resp, err := http.Get("http://localhost:11000/api/config")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	cfg := config.Config{}
	err = json.NewDecoder(resp.Body).Decode(&cfg)
	require.NoError(t, err)

	assert.False(t, cfg.DryRunMode)
	assert.Equal(t, 11000, cfg.Port)
	assert.Equal(t, "127.0.0.1", cfg.Scaler.Nomad.ServerAddr)
	assert.Equal(t, time.Second*150, cfg.CapacityPlanner.DownScaleCooldownPeriod)
	assert.Equal(t, time.Second*160, cfg.CapacityPlanner.UpScaleCooldownPeriod)
	assert.Equal(t, 2, int(cfg.ScaleAlertAggregator.NoAlertScaleDamping))
	assert.Equal(t, 30, int(cfg.ScaleAlertAggregator.UpScaleThreshold))
	assert.Equal(t, -40, int(cfg.ScaleAlertAggregator.DownScaleThreshold))
	assert.Equal(t, time.Second*15, cfg.ScaleAlertAggregator.EvaluationCycle)
	assert.Equal(t, 25, int(cfg.ScaleAlertAggregator.EvaluationPeriodFactor))
	assert.Equal(t, time.Second*260, cfg.ScaleAlertAggregator.CleanupCycle)
}

func Test_EndPoints_Metrics(t *testing.T) {
	resp, err := http.Get("http://localhost:11000/metrics")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_EndPoints_Health(t *testing.T) {
	resp, err := http.Get("http://localhost:11000/health")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_EndPoints_Build(t *testing.T) {
	resp, err := http.Get("http://localhost:11000/api/build")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
