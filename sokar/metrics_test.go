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
	metrics := Metrics{
		scaleEventsTotal:   mScaleEventsTotal,
		failedScalingTotal: mFailedScalingTotal,
		plannedCount:       mPlannedCount,
	}
	mocks := MetricsMocks{
		scaleEventsTotal:   mScaleEventsTotal,
		failedScalingTotal: mFailedScalingTotal,
		plannedCount:       mPlannedCount,
	}
	return metrics, mocks
}

func Test_NewMetrics(t *testing.T) {
	metrics := NewMetrics()
	assert.NotNil(t, metrics.scaleEventsTotal)
	assert.NotNil(t, metrics.failedScalingTotal)
	assert.NotNil(t, metrics.plannedCount)
}

//func Test_UpdateAlertMetrics(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	metrics, mocks := NewMockedMetrics(mockCtrl)
//
//	gaugeUP := mock_metrics.NewMockGauge(mockCtrl)
//	gaugeUP.EXPECT().Set(float64(1))
//	gaugeDOWN := mock_metrics.NewMockGauge(mockCtrl)
//	gaugeDOWN.EXPECT().Set(float64(2))
//	gaugeNeutral := mock_metrics.NewMockGauge(mockCtrl)
//	gaugeNeutral.EXPECT().Set(float64(3))
//
//	gomock.InOrder(
//		mocks.alerts.EXPECT().WithLabelValues("up").Return(gaugeUP),
//		mocks.alerts.EXPECT().WithLabelValues("down").Return(gaugeDOWN),
//		mocks.alerts.EXPECT().WithLabelValues("neutral").Return(gaugeNeutral),
//	)
//
//	scap := NewScaleAlertPool()
//	scap.entries[1] = ScaleAlertPoolEntry{weight: 1}
//	scap.entries[2] = ScaleAlertPoolEntry{weight: -1}
//	scap.entries[3] = ScaleAlertPoolEntry{weight: -100}
//	scap.entries[4] = ScaleAlertPoolEntry{weight: 0}
//	scap.entries[5] = ScaleAlertPoolEntry{weight: 0}
//	scap.entries[6] = ScaleAlertPoolEntry{weight: 0}
//	updateAlertMetrics(&scap, &metrics)
//}
//
