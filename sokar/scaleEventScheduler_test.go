package sokar

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_metrics "github.com/thomasobenaus/sokar/test/mocks/metrics"
	mock_sokar "github.com/thomasobenaus/sokar/test/mocks/sokar"
)

func Test_ShouldFireScaleEvent(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)
	metrics, mocks := NewMockedMetrics(mockCtrl)
	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)
	minScale := uint(1)
	maxScale := uint(10)

	// WHEN
	scheduleIF.EXPECT().IsActiveAt(gomock.Any(), gomock.Any()).Return(true)
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale, maxScale, nil)
	scheduledScaleBoundsMin := mock_metrics.NewMockGauge(mockCtrl)
	scheduledScaleBoundsMin.EXPECT().Set(float64(minScale))
	scheduledScaleBoundsMax := mock_metrics.NewMockGauge(mockCtrl)
	scheduledScaleBoundsMax.EXPECT().Set(float64(maxScale))
	mocks.scheduledScaleBounds.EXPECT().WithLabelValues("min").Return(scheduledScaleBoundsMin)
	mocks.scheduledScaleBounds.EXPECT().WithLabelValues("max").Return(scheduledScaleBoundsMax)
	result := sokar.shouldFireScaleEvent(time.Now())

	// THEN
	assert.True(t, result)
}

func Test_ShouldNotFireScaleEvent(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)
	metrics, mocks := NewMockedMetrics(mockCtrl)
	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)
	minScale := uint(1)
	maxScale := uint(10)
	err1 := fmt.Errorf("No scale schedule entry")

	// WHEN
	scheduleIF.EXPECT().IsActiveAt(gomock.Any(), gomock.Any()).Return(false)
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale, maxScale, err1)
	scheduledScaleBoundsMin := mock_metrics.NewMockGauge(mockCtrl)
	scheduledScaleBoundsMin.EXPECT().Set(float64(0))
	scheduledScaleBoundsMax := mock_metrics.NewMockGauge(mockCtrl)
	scheduledScaleBoundsMax.EXPECT().Set(float64(0))
	mocks.scheduledScaleBounds.EXPECT().WithLabelValues("min").Return(scheduledScaleBoundsMin)
	mocks.scheduledScaleBounds.EXPECT().WithLabelValues("max").Return(scheduledScaleBoundsMax)
	result := sokar.shouldFireScaleEvent(time.Now())

	// THEN
	assert.False(t, result)
}
