package scaler

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_metrics "github.com/thomasobenaus/sokar/test/metrics"
	mock_scaler "github.com/thomasobenaus/sokar/test/scaler"
)

func TestScale_ScalingObjectDead(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	cfg := Config{Name: sObjName, WatcherInterval: time.Second * 5}
	scaler, err := cfg.New(scaTgt, ScalingObject{Name: sObjName}, metrics)
	require.NoError(t, err)

	// dead scalingObject - error
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, fmt.Errorf("internal error"))
	result := scaler.scale(2, 0, false)
	assert.Equal(t, scaleFailed, result.state)

	// dead scalingObject
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(true, nil)
	result = scaler.scale(2, 0, false)
	assert.Equal(t, scaleIgnored, result.state)
}

func TestScale_Up(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	cfg := Config{Name: sObjName, MinCount: 1, MaxCount: 5, WatcherInterval: time.Second * 5}
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0)).Times(2)
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge).Times(2)

	// scale up
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(0), uint(2)).Return(nil)
	result := scaler.scale(2, 0, false)
	assert.Equal(t, uint(2), scaler.desiredScale.value)
	assert.NotEqual(t, scaleFailed, result.state)

	// scale up - max hit
	polViolatedCounter := mock_metrics.NewMockCounter(mockCtrl)
	polViolatedCounter.EXPECT().Inc()
	mocks.scalingPolicyViolated.EXPECT().WithLabelValues("max").Return(polViolatedCounter)
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(0), uint(5)).Return(nil)
	result = scaler.scale(6, 0, false)
	assert.Equal(t, uint(5), scaler.desiredScale.value)
	assert.NotEqual(t, scaleFailed, result.state)
}

func TestScale_Down(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	cfg := Config{Name: sObjName, MinCount: 1, MaxCount: 5, WatcherInterval: time.Second * 5}
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0)).Times(2)
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("DOWN").Return(plannedButSkippedGauge).Times(2)

	// scale down
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(4), uint(1)).Return(nil)
	result := scaler.scale(1, 4, false)
	assert.Equal(t, uint(1), scaler.desiredScale.value)
	assert.NotEqual(t, scaleFailed, result.state)

	// scale up - min hit
	polViolatedCounter := mock_metrics.NewMockCounter(mockCtrl)
	polViolatedCounter.EXPECT().Inc()
	mocks.scalingPolicyViolated.EXPECT().WithLabelValues("min").Return(polViolatedCounter)
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(2), uint(1)).Return(nil)
	result = scaler.scale(0, 2, false)
	assert.Equal(t, uint(1), scaler.desiredScale.value)
	assert.NotEqual(t, scaleFailed, result.state)
}

func TestScale_NoScale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	cfg := Config{Name: sObjName, MinCount: 1, MaxCount: 5, WatcherInterval: time.Second * 5}
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	// scale down
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	result := scaler.scale(2, 2, false)
	assert.False(t, scaler.desiredScale.isKnown)
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

	sObjName := "any"
	cfg := Config{Name: sObjName, MinCount: 1, MaxCount: 5, WatcherInterval: time.Second * 5}
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(1))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge)

	// scale up
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	result := scaler.scale(2, 0, true)
	assert.NotEqual(t, scaleFailed, result.state)
}

func TestScale_DownDryRun(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	cfg := Config{Name: sObjName, MinCount: 1, MaxCount: 5, WatcherInterval: time.Second * 5}
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(1))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("DOWN").Return(plannedButSkippedGauge)

	// scale down
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	result := scaler.scale(1, 4, true)
	assert.NotEqual(t, scaleFailed, result.state)
}
