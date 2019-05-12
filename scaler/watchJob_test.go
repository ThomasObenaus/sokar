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

func TestCountMeetsExpectations(t *testing.T) {

	desired := uint(1)
	asExpected, expected := countMeetsExpectations(1, 1, 2, &desired)
	assert.True(t, asExpected)
	assert.Equal(t, uint(1), expected)

	// min violated
	asExpected, expected = countMeetsExpectations(0, 1, 2, &desired)
	assert.False(t, asExpected)
	assert.Equal(t, uint(1), expected)

	// max violated
	asExpected, expected = countMeetsExpectations(3, 1, 2, &desired)
	assert.False(t, asExpected)
	assert.Equal(t, uint(2), expected)

	// no desired
	asExpected, expected = countMeetsExpectations(1, 1, 2, nil)
	assert.True(t, asExpected)
	assert.Equal(t, uint(1), expected, "Expected to stay at current count since desired is nil.")

	// desired over max
	desired = uint(10)
	asExpected, expected = countMeetsExpectations(3, 1, 2, &desired)
	assert.False(t, asExpected)
	assert.Equal(t, uint(2), expected)

	// desired under min
	desired = uint(0)
	asExpected, expected = countMeetsExpectations(3, 1, 2, &desired)
	assert.False(t, asExpected)
	assert.Equal(t, uint(2), expected)
}

func TestEnsureJobCount_NoScale(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg := Config{JobName: "any", MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// Error in scaling target
	scaTgt.EXPECT().GetJobCount("any").Return(uint(0), fmt.Errorf("Unable to determine job count"))
	err = scaler.ensureJobCount()
	assert.Error(t, err)

	// No scale
	desired := uint(1)
	scaler.desiredScale = &desired
	scaTgt.EXPECT().GetJobCount("any").Return(uint(1), nil)
	err = scaler.ensureJobCount()
	assert.NoError(t, err)

}

func TestEnsureJobCount_MinViolated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg := Config{JobName: "any", MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// Scale - min violated
	desired := uint(1)
	scaler.desiredScale = &desired
	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc()
	mocks.scalingTicketCount.EXPECT().WithLabelValues("added").Return(scalingTicketCounter)
	scaTgt.EXPECT().GetJobCount("any").Return(uint(0), nil)
	err = scaler.ensureJobCount()
	assert.NoError(t, err)

}

func TestEnsureJobCount_MaxViolated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg := Config{JobName: "any", MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// Scale - max violated
	desired := uint(1)
	scaler.desiredScale = &desired
	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc()
	mocks.scalingTicketCount.EXPECT().WithLabelValues("added").Return(scalingTicketCounter)
	scaTgt.EXPECT().GetJobCount("any").Return(uint(10), nil)
	err = scaler.ensureJobCount()
	assert.NoError(t, err)
}
