package scaler

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/test/metrics"
	"github.com/thomasobenaus/sokar/test/scaler"
)

func TestScale_JobDead(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)

	// dead job - error
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, fmt.Errorf("internal error"))
	result := scaler.scale(2, 0, false)
	assert.Equal(t, scaleFailed, result.state)

	// dead job
	scaTgt.EXPECT().IsJobDead(jobname).Return(true, nil)
	result = scaler.scale(2, 0, false)
	assert.Equal(t, scaleIgnored, result.state)
}

func TestScale_Up(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0)).Times(2)
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge).Times(2)

	// scale up
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(2)).Return(nil)
	result := scaler.scale(2, 0, false)
	assert.NotEqual(t, scaleFailed, result.state)

	// scale up - max hit
	polViolatedCounter := mock_metrics.NewMockCounter(mockCtrl)
	polViolatedCounter.EXPECT().Inc()
	mocks.scalingPolicyViolated.EXPECT().WithLabelValues("max").Return(polViolatedCounter)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(5)).Return(nil)
	result = scaler.scale(6, 0, false)
	assert.NotEqual(t, scaleFailed, result.state)
}

func TestScale_Down(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0)).Times(2)
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("DOWN").Return(plannedButSkippedGauge).Times(2)

	// scale down
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(1)).Return(nil)
	result := scaler.scale(1, 4, false)
	assert.NotEqual(t, scaleFailed, result.state)

	// scale up - min hit
	polViolatedCounter := mock_metrics.NewMockCounter(mockCtrl)
	polViolatedCounter.EXPECT().Inc()
	mocks.scalingPolicyViolated.EXPECT().WithLabelValues("min").Return(polViolatedCounter)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(1)).Return(nil)
	result = scaler.scale(0, 2, false)
	assert.NotEqual(t, scaleFailed, result.state)
}

func TestScale_NoScale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)

	// scale down
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	result := scaler.scale(2, 2, false)
	assert.NotEqual(t, scaleFailed, result.state)
}

func TestScaleBy_CheckScalingPolicy(t *testing.T) {

	chk := checkScalingPolicy(0, 0, 0)
	assert.Equal(t, uint(0), chk.validCount)
	assert.Equal(t, uint(0), chk.desiredCount)
	assert.False(t, chk.minPolicyViolated)
	assert.False(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(1, 0, 0)
	assert.Equal(t, uint(0), chk.validCount)
	assert.Equal(t, uint(1), chk.desiredCount)
	assert.False(t, chk.minPolicyViolated)
	assert.True(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(1, 2, 3)
	assert.Equal(t, uint(2), chk.validCount)
	assert.Equal(t, uint(1), chk.desiredCount)
	assert.True(t, chk.minPolicyViolated)
	assert.False(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(3, 1, 4)
	assert.Equal(t, uint(3), chk.validCount)
	assert.Equal(t, uint(3), chk.desiredCount)
	assert.False(t, chk.minPolicyViolated)
	assert.False(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(3, 1, 2)
	assert.Equal(t, uint(2), chk.validCount)
	assert.Equal(t, uint(3), chk.desiredCount)
	assert.False(t, chk.minPolicyViolated)
	assert.True(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(2, 3, 1)
	assert.Equal(t, uint(1), chk.validCount)
	assert.Equal(t, uint(2), chk.desiredCount)
	assert.True(t, chk.minPolicyViolated)
	assert.True(t, chk.maxPolicyViolated)
}

func TestScaleBy_trueIfNil(t *testing.T) {
	_, ok := trueIfNil(nil)
	assert.True(t, ok)

	scaler := &Scaler{}
	_, ok = trueIfNil(scaler)
	assert.False(t, ok)
}

func TestScale_UpDryRun(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(1))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge)

	// scale up
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	result := scaler.scale(2, 0, true)
	assert.NotEqual(t, scaleFailed, result.state)
}

func TestScale_DownDryRun(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(1))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("DOWN").Return(plannedButSkippedGauge)

	// scale down
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	result := scaler.scale(1, 4, true)
	assert.NotEqual(t, scaleFailed, result.state)
}
