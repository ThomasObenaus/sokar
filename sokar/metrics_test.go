package sokar

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/metrics"
)

type MetricsMocks struct {
	scaleEventsTotal   *mock_metrics.MockCounter
	failedScalingTotal *mock_metrics.MockCounter
	plannedCount       *mock_metrics.MockGauge
	currentCount       *mock_metrics.MockGauge
	scaleFactor        *mock_metrics.MockGauge
}

// NewMockedMetrics creates and returns mocked metrics that can be used
// for unit-testing.
// Example:
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()
// 		metrics, mocks := NewMockedMetrics(mockCtrl)
// 		mocks.scaleCounter.EXPECT().Set(10)
// use metrics...
func NewMockedMetrics(mockCtrl *gomock.Controller) (Metrics, MetricsMocks) {
	mScaleEventsTotal := mock_metrics.NewMockCounter(mockCtrl)
	mFailedScalingTotal := mock_metrics.NewMockCounter(mockCtrl)
	mPlannedCount := mock_metrics.NewMockGauge(mockCtrl)
	mCurrentCount := mock_metrics.NewMockGauge(mockCtrl)
	mScaleFactor := mock_metrics.NewMockGauge(mockCtrl)
	metrics := Metrics{
		scaleEventsTotal:   mScaleEventsTotal,
		failedScalingTotal: mFailedScalingTotal,
		plannedCount:       mPlannedCount,
		currentCount:       mCurrentCount,
		scaleFactor:        mScaleFactor,
	}
	mocks := MetricsMocks{
		scaleEventsTotal:   mScaleEventsTotal,
		failedScalingTotal: mFailedScalingTotal,
		plannedCount:       mPlannedCount,
		currentCount:       mCurrentCount,
		scaleFactor:        mScaleFactor,
	}
	return metrics, mocks
}

func Test_NewMetrics(t *testing.T) {
	metrics := NewMetrics()
	assert.NotNil(t, metrics.scaleEventsTotal)
	assert.NotNil(t, metrics.failedScalingTotal)
	assert.NotNil(t, metrics.plannedCount)
	assert.NotNil(t, metrics.currentCount)
	assert.NotNil(t, metrics.scaleFactor)
}
