package scaler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/metrics"
)

type MetricsMocks struct {
	scalingPolicyViolated *mock_metrics.MockCounterVec
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
	mScalingPolicyViolated := mock_metrics.NewMockCounterVec(mockCtrl)
	metrics := Metrics{
		scalingPolicyViolated: mScalingPolicyViolated,
	}
	mocks := MetricsMocks{
		scalingPolicyViolated: mScalingPolicyViolated,
	}
	return metrics, mocks
}

func Test_NewMetrics(t *testing.T) {
	metrics := NewMetrics()
	assert.NotNil(t, metrics.scalingPolicyViolated)
}
