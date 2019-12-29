package capacityPlanner

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_metrics "github.com/thomasobenaus/sokar/test/metrics"
)

type MetricsMocks struct {
	scheduledScaleBounds *mock_metrics.MockGaugeVec
	scaleAdjustments     *mock_metrics.MockGaugeVec
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
	mScaleAdjustments := mock_metrics.NewMockGaugeVec(mockCtrl)
	metrics := Metrics{
		scheduledScaleBounds: mScheduledScaleBounds,
		scaleAdjustments:     mScaleAdjustments,
	}
	mocks := MetricsMocks{
		scheduledScaleBounds: mScheduledScaleBounds,
		scaleAdjustments:     mScaleAdjustments,
	}
	return metrics, mocks
}

func Test_NewMetrics(t *testing.T) {
	metrics := NewMetrics()
	assert.NotNil(t, metrics.scheduledScaleBounds)
	assert.NotNil(t, metrics.scaleAdjustments)
}
