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

func TestCountMeetsExpectations(t *testing.T) {

	desired := optionalValue{value: 1, isKnown: true}
	asExpected, expected := countMeetsExpectations(1, 1, 2, desired)
	assert.True(t, asExpected)
	assert.Equal(t, uint(1), expected)

	// min violated
	asExpected, expected = countMeetsExpectations(0, 1, 2, desired)
	assert.False(t, asExpected)
	assert.Equal(t, uint(1), expected)

	// max violated
	asExpected, expected = countMeetsExpectations(3, 1, 2, desired)
	assert.False(t, asExpected)
	assert.Equal(t, uint(2), expected)

	// no desired
	asExpected, expected = countMeetsExpectations(1, 1, 2, optionalValue{isKnown: false})
	assert.True(t, asExpected)
	assert.Equal(t, uint(1), expected, "Expected to stay at current count since desired is nil.")

	// desired over max
	desired.setValue(10)
	asExpected, expected = countMeetsExpectations(3, 1, 2, desired)
	assert.False(t, asExpected)
	assert.Equal(t, uint(2), expected)

	// desired under min
	desired.setValue(0)
	asExpected, expected = countMeetsExpectations(3, 1, 2, desired)
	assert.False(t, asExpected)
	assert.Equal(t, uint(2), expected)
}

func TestEnsureScalingObjectCount_NoScale(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	cfg := Config{Name: sObjName, MinCount: 1, MaxCount: 5, WatcherInterval: time.Second * 5}
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, sObj, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// Error in scaling target
	scaTgt.EXPECT().GetScalingObjectCount(sObjName).Return(uint(0), fmt.Errorf("Unable to determine scalingObject count"))
	err = scaler.ensureScalingObjectCount()
	assert.Error(t, err)

	// No scale
	scaler.desiredScale.setValue(1)
	scaTgt.EXPECT().GetScalingObjectCount(sObjName).Return(uint(1), nil)
	err = scaler.ensureScalingObjectCount()
	assert.NoError(t, err)

}

func TestEnsureScalingObjectCount_MinViolated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	cfg := Config{Name: sObjName, MinCount: 1, MaxCount: 5, WatcherInterval: time.Second * 5}
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, sObj, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// Scale - min violated
	scaler.desiredScale.setValue(1)
	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc()
	mocks.scalingTicketCount.EXPECT().WithLabelValues("added").Return(scalingTicketCounter)
	scaTgt.EXPECT().GetScalingObjectCount("any").Return(uint(0), nil)
	err = scaler.ensureScalingObjectCount()
	assert.NoError(t, err)
}

func TestEnsureScalingObjectCount_MaxViolated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObjName := "any"
	cfg := Config{Name: sObjName, MinCount: 1, MaxCount: 5, WatcherInterval: time.Second * 5}
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, sObj, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// Scale - max violated
	scaler.desiredScale.setValue(1)
	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc()
	mocks.scalingTicketCount.EXPECT().WithLabelValues("added").Return(scalingTicketCounter)
	scaTgt.EXPECT().GetScalingObjectCount(sObjName).Return(uint(10), nil)
	err = scaler.ensureScalingObjectCount()
	assert.NoError(t, err)
}
