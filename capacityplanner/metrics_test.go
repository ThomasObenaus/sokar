package capacityplanner

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_metrics "github.com/thomasobenaus/sokar/test/mocks/metrics"
)

type MetricsMocks struct {
	scaleAdjustments *mock_metrics.MockGaugeVec
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
	mScaleAdjustments := mock_metrics.NewMockGaugeVec(mockCtrl)
	metrics := Metrics{
		scaleAdjustments: mScaleAdjustments,
	}
	mocks := MetricsMocks{
		scaleAdjustments: mScaleAdjustments,
	}
	return metrics, mocks
}

func Test_NewMetrics(t *testing.T) {
	metrics := NewMetrics()
	assert.NotNil(t, metrics.scaleAdjustments)
}
