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

func Test_ExecuteScale_NoDryRun(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sObjName := "any"
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	metrics, mocks := NewMockedMetrics(mockCtrl)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0)).Times(3)
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge).Times(3)

	sObj := ScalingObject{Name: sObjName, MinCount: 0, MaxCount: 2}
	scaler, err := New(scaTgt, sObj, metrics, DryRunMode(false))
	require.NoError(t, err)

	// error
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(0), uint(2), uint(1), uint(1)).Return(fmt.Errorf("err"))
	result := scaler.executeScale(1, 1, false)
	assert.Equal(t, scaleFailed, result.state)
	assert.Equal(t, uint(1), result.newCount)

	// no scale, success
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(0), uint(2), uint(1), uint(1)).Return(nil)
	result = scaler.executeScale(1, 1, false)
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(1), result.newCount)

	// scale up, success
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(0), uint(2), uint(1), uint(2)).Return(nil)
	result = scaler.executeScale(1, 2, false)
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(2), result.newCount)

	// scale down, success
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(0), uint(2), uint(2), uint(1)).Return(nil)
	plannedButSkippedGauge.EXPECT().Set(float64(0))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("DOWN").Return(plannedButSkippedGauge)
	result = scaler.executeScale(2, 1, false)
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(1), result.newCount)
}

func Test_ExecuteScale_DryRun(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sObjName := "any"
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	metrics, mocks := NewMockedMetrics(mockCtrl)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(1))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge)

	sObj := ScalingObject{Name: sObjName, MinCount: 0, MaxCount: 2}
	scaler, err := New(scaTgt, sObj, metrics, DryRunMode(true))
	require.NoError(t, err)

	// no scale, dry run
	result := scaler.executeScale(1, 2, false)
	assert.Equal(t, scaleIgnored, result.state)
	assert.Equal(t, uint(1), result.newCount)

	// scale up, dry run but force
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(0), uint(2), uint(1), uint(2)).Return(nil)
	plannedButSkippedGauge.EXPECT().Set(float64(0))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge)
	result = scaler.executeScale(1, 2, true)
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(2), result.newCount)
}

func TestScale_ScalingObjectDead(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	scaler, err := New(scaTgt, ScalingObject{Name: sObjName}, metrics)
	require.NoError(t, err)

	// dead scalingObject - error
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, fmt.Errorf("internal error"))
	result := scaler.scale(2, 1, false)
	assert.Equal(t, scaleFailed, result.state)
	assert.Equal(t, uint(1), result.newCount)

	// dead scalingObject
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(true, nil)
	result = scaler.scale(2, 1, false)
	assert.Equal(t, scaleIgnored, result.state)
	assert.Equal(t, uint(1), result.newCount)
}

func TestScale_Up(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0)).Times(2)
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge).Times(2)

	// scale up
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(0), uint(2)).Return(nil)
	result := scaler.scale(2, 0, false)
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(2), result.newCount)

	// scale up - max hit
	polViolatedCounter := mock_metrics.NewMockCounter(mockCtrl)
	polViolatedCounter.EXPECT().Inc()
	mocks.scalingPolicyViolated.EXPECT().WithLabelValues("max").Return(polViolatedCounter)
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(0), uint(5)).Return(nil)
	result = scaler.scale(6, 0, false)
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(5), result.newCount)
}

func TestScale_Down(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0)).Times(2)
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("DOWN").Return(plannedButSkippedGauge).Times(2)

	// scale down
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(4), uint(1)).Return(nil)
	result := scaler.scale(1, 4, false)
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(1), result.newCount)

	// scale up - min hit
	polViolatedCounter := mock_metrics.NewMockCounter(mockCtrl)
	polViolatedCounter.EXPECT().Inc()
	mocks.scalingPolicyViolated.EXPECT().WithLabelValues("min").Return(polViolatedCounter)
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(2), uint(1)).Return(nil)
	result = scaler.scale(0, 2, false)
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(1), result.newCount)
}

func TestScale_NoScale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	result := scaler.scale(2, 2, false)
	assert.Equal(t, scaleIgnored, result.state)
	assert.Equal(t, uint(2), result.newCount)
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

func TestScale_UpDryRun(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := New(scaTgt, sObj, metrics, DryRunMode(true))
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(1))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge)

	// scale up
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	result := scaler.scale(2, 1, false)
	assert.Equal(t, scaleIgnored, result.state)
	assert.Equal(t, uint(1), result.newCount)
	oneDayAgo := time.Now().Add(time.Hour * -24)
	assert.WithinDuration(t, oneDayAgo, scaler.lastScaleAction, time.Second*1)
}

func TestScale_DryRunForce(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := New(scaTgt, sObj, metrics, DryRunMode(true))
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(1), uint(2)).Return(nil)

	// scale up
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	result := scaler.scale(2, 1, true)
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(2), result.newCount)
	fiveMsAgo := time.Now().Add(time.Millisecond * -5)
	assert.WithinDuration(t, fiveMsAgo, scaler.lastScaleAction, time.Second*1)
}

func TestScale_DownDryRun(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := New(scaTgt, sObj, metrics, DryRunMode(true))
	require.NoError(t, err)

	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(1))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("DOWN").Return(plannedButSkippedGauge)

	// scale down
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)
	result := scaler.scale(1, 4, false)
	assert.Equal(t, scaleIgnored, result.state)
	assert.Equal(t, uint(4), result.newCount)
	oneDayAgo := time.Now().Add(time.Hour * -24)
	assert.WithinDuration(t, oneDayAgo, scaler.lastScaleAction, time.Second*1)
}

func Test_IsScalePermitted(t *testing.T) {
	assert.True(t, isScalePermitted(true, true))
	assert.True(t, isScalePermitted(false, true))
	assert.True(t, isScalePermitted(false, false))
	assert.False(t, isScalePermitted(true, false))
}

func Test_ShouldUpdateLastScaleActionOnScale(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	sObjName := "any"
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	// WHEN - expecting
	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0))
	mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("DOWN").Return(plannedButSkippedGauge)
	scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(5), uint(4), uint(2)).Return(nil)
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)

	// WHEN - scale down
	result := scaler.scale(2, 4, false)

	// THEN
	assert.Equal(t, scaleDone, result.state)
	assert.Equal(t, uint(2), result.newCount)
	assert.WithinDuration(t, time.Now(), scaler.lastScaleAction, time.Second*1)
	//oneDayAgo := time.Now().Add(time.Hour * -24)
}

func Test_ShouldNotUpdateLastScaleActionIfNoScaleIsNeeded(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	sObjName := "any"
	sObj := ScalingObject{Name: sObjName, MinCount: 2, MaxCount: 5}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)

	// WHEN - expecting
	polViolatedCounter := mock_metrics.NewMockCounter(mockCtrl)
	polViolatedCounter.EXPECT().Inc()
	mocks.scalingPolicyViolated.EXPECT().WithLabelValues("min").Return(polViolatedCounter)
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil)

	// WHEN - scale
	result := scaler.scale(1, 2, false)

	// THEN
	assert.Equal(t, scaleIgnored, result.state)
	assert.Equal(t, uint(2), result.newCount)
	oneDayAgo := time.Now().Add(time.Hour * -24)
	assert.WithinDuration(t, oneDayAgo, scaler.lastScaleAction, time.Second*1)
}
