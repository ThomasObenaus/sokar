package scaler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/mocks/metrics"
)

type MetricsMocks struct {
	scalingPolicyViolated        *mock_metrics.MockCounterVec
	scalingTicketCount           *mock_metrics.MockCounterVec
	scaleResultCounter           *mock_metrics.MockCounterVec
	scalingDurationSeconds       *mock_metrics.MockHistogram
	plannedButSkippedScalingOpen *mock_metrics.MockGaugeVec
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
	mScalingTicketCount := mock_metrics.NewMockCounterVec(mockCtrl)
	mScaleResultCounter := mock_metrics.NewMockCounterVec(mockCtrl)
	mScalingDurationSeconds := mock_metrics.NewMockHistogram(mockCtrl)
	mPlannedButSkippedScalingOpen := mock_metrics.NewMockGaugeVec(mockCtrl)
	metrics := Metrics{
		scalingPolicyViolated:        mScalingPolicyViolated,
		scalingTicketCount:           mScalingTicketCount,
		scaleResultCounter:           mScaleResultCounter,
		scalingDurationSeconds:       mScalingDurationSeconds,
		plannedButSkippedScalingOpen: mPlannedButSkippedScalingOpen,
	}
	mocks := MetricsMocks{
		scalingPolicyViolated:        mScalingPolicyViolated,
		scalingTicketCount:           mScalingTicketCount,
		scaleResultCounter:           mScaleResultCounter,
		scalingDurationSeconds:       mScalingDurationSeconds,
		plannedButSkippedScalingOpen: mPlannedButSkippedScalingOpen,
	}
	return metrics, mocks
}

func Test_NewMetrics(t *testing.T) {
	metrics := NewMetrics()
	assert.NotNil(t, metrics.scalingPolicyViolated)
	assert.NotNil(t, metrics.scalingTicketCount)
	assert.NotNil(t, metrics.scaleResultCounter)
	assert.NotNil(t, metrics.scalingDurationSeconds)
	assert.NotNil(t, metrics.plannedButSkippedScalingOpen)
}
