package sokar

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_metrics "github.com/thomasobenaus/sokar/test/metrics"
)

type MetricsMocks struct {
	scheduledScaleBounds              *mock_metrics.MockGaugeVec
	scaleEventsTotal                  *mock_metrics.MockCounter
	failedScalingTotal                *mock_metrics.MockCounter
	skippedScalingDuringCooldownTotal *mock_metrics.MockCounter
	preScaleJobCount                  *mock_metrics.MockGauge
	plannedJobCount                   *mock_metrics.MockGauge
	scaleFactor                       *mock_metrics.MockGauge
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
	mScheduledScaleBounds := mock_metrics.NewMockGaugeVec(mockCtrl)
	mScaleEventsTotal := mock_metrics.NewMockCounter(mockCtrl)
	mFailedScalingTotal := mock_metrics.NewMockCounter(mockCtrl)
	mSkippedScalingDuringCooldownTotal := mock_metrics.NewMockCounter(mockCtrl)
	mPlannedJobCount := mock_metrics.NewMockGauge(mockCtrl)
	mPreScaleJobCount := mock_metrics.NewMockGauge(mockCtrl)
	mScaleFactor := mock_metrics.NewMockGauge(mockCtrl)
	metrics := Metrics{
		scheduledScaleBounds:              mScheduledScaleBounds,
		scaleEventsTotal:                  mScaleEventsTotal,
		failedScalingTotal:                mFailedScalingTotal,
		skippedScalingDuringCooldownTotal: mSkippedScalingDuringCooldownTotal,
		plannedJobCount:                   mPlannedJobCount,
		preScaleJobCount:                  mPreScaleJobCount,
		scaleFactor:                       mScaleFactor,
	}
	mocks := MetricsMocks{
		scheduledScaleBounds:              mScheduledScaleBounds,
		scaleEventsTotal:                  mScaleEventsTotal,
		failedScalingTotal:                mFailedScalingTotal,
		skippedScalingDuringCooldownTotal: mSkippedScalingDuringCooldownTotal,
		preScaleJobCount:                  mPreScaleJobCount,
		plannedJobCount:                   mPlannedJobCount,
		scaleFactor:                       mScaleFactor,
	}
	return metrics, mocks
}

func Test_NewMetrics(t *testing.T) {
	metrics := NewMetrics()
	assert.NotNil(t, metrics.scheduledScaleBounds)
	assert.NotNil(t, metrics.scaleEventsTotal)
	assert.NotNil(t, metrics.failedScalingTotal)
	assert.NotNil(t, metrics.skippedScalingDuringCooldownTotal)
	assert.NotNil(t, metrics.preScaleJobCount)
	assert.NotNil(t, metrics.plannedJobCount)
	assert.NotNil(t, metrics.scaleFactor)
}
