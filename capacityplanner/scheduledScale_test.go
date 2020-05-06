package capacityplanner

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_capacityplanner "github.com/thomasobenaus/sokar/test/mocks/capacityplanner"
	mock_metrics "github.com/thomasobenaus/sokar/test/mocks/metrics"
)

func Test_FitIntoScaleRangeShouldNotAdjustIfItIsInBounds(t *testing.T) {
	// GIVEN
	plannedScale1 := uint(10)
	plannedScale2 := uint(1)
	plannedScale3 := uint(5)
	minScale := uint(1)
	maxScale := uint(10)

	// WHEN
	adjustedScale1 := fitIntoScaleRange(plannedScale1, minScale, maxScale)
	adjustedScale2 := fitIntoScaleRange(plannedScale2, minScale, maxScale)
	adjustedScale3 := fitIntoScaleRange(plannedScale3, minScale, maxScale)

	// THEN
	assert.Equal(t, uint(10), adjustedScale1)
	assert.Equal(t, uint(1), adjustedScale2)
	assert.Equal(t, uint(5), adjustedScale3)
}

func Test_FitIntoScaleRangeShouldAdjustIfItIsNotInBounds(t *testing.T) {
	// GIVEN
	plannedScale1 := uint(0)
	plannedScale2 := uint(11)
	plannedScale3 := uint(20)
	minScale := uint(1)
	maxScale := uint(10)

	// WHEN
	adjustedScale1 := fitIntoScaleRange(plannedScale1, minScale, maxScale)
	adjustedScale2 := fitIntoScaleRange(plannedScale2, minScale, maxScale)
	adjustedScale3 := fitIntoScaleRange(plannedScale3, minScale, maxScale)

	// THEN
	assert.Equal(t, uint(1), adjustedScale1)
	assert.Equal(t, uint(10), adjustedScale2)
	assert.Equal(t, uint(10), adjustedScale3)
}

func Test_ShouldAdjustScale(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	scheduleIF := mock_capacityplanner.NewMockScaleSchedule(mockCtrl)
	metrics, mocks := NewMockedMetrics(mockCtrl)

	capa, err := New(metrics, Schedule(scheduleIF))
	minScale := uint(1)
	maxScale := uint(10)
	plannedScale1 := uint(0)

	// WHEN
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale, maxScale, nil)

	scaleAdjustmentsPlanned := mock_metrics.NewMockGauge(mockCtrl)
	scaleAdjustmentsPlanned.EXPECT().Set(float64(0))
	mocks.scaleAdjustments.EXPECT().WithLabelValues("planned").Return(scaleAdjustmentsPlanned)
	scaleAdjustmentsAdjusted := mock_metrics.NewMockGauge(mockCtrl)
	scaleAdjustmentsAdjusted.EXPECT().Set(float64(1))
	mocks.scaleAdjustments.EXPECT().WithLabelValues("adjusted").Return(scaleAdjustmentsAdjusted)

	replannedScale1 := capa.adjustPlanAccordingToSchedule(plannedScale1, time.Now())

	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, capa)
	assert.Equal(t, uint(1), replannedScale1)

	// GIVEN
	plannedScale2 := uint(11)

	// WHEN
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale, maxScale, nil)
	scaleAdjustmentsPlanned.EXPECT().Set(float64(11))
	mocks.scaleAdjustments.EXPECT().WithLabelValues("planned").Return(scaleAdjustmentsPlanned)
	scaleAdjustmentsAdjusted.EXPECT().Set(float64(10))
	mocks.scaleAdjustments.EXPECT().WithLabelValues("adjusted").Return(scaleAdjustmentsAdjusted)

	replannedScale2 := capa.adjustPlanAccordingToSchedule(plannedScale2, time.Now())

	// THEN
	assert.Equal(t, uint(10), replannedScale2)
}

func Test_ShouldNotAdjustScale(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	scheduleIF := mock_capacityplanner.NewMockScaleSchedule(mockCtrl)
	metrics, mocks := NewMockedMetrics(mockCtrl)
	capa, err := New(metrics, Schedule(scheduleIF))
	plannedScale := uint(5)
	minScale := uint(1)
	maxScale := uint(10)
	err1 := fmt.Errorf("some err")

	// WHEN
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale, maxScale, err1)
	scaleAdjustmentsPlanned := mock_metrics.NewMockGauge(mockCtrl)
	scaleAdjustmentsPlanned.EXPECT().Set(float64(5))
	mocks.scaleAdjustments.EXPECT().WithLabelValues("planned").Return(scaleAdjustmentsPlanned)
	scaleAdjustmentsAdjusted := mock_metrics.NewMockGauge(mockCtrl)
	scaleAdjustmentsAdjusted.EXPECT().Set(float64(5))
	mocks.scaleAdjustments.EXPECT().WithLabelValues("adjusted").Return(scaleAdjustmentsAdjusted)

	replannedScale1 := capa.adjustPlanAccordingToSchedule(plannedScale, time.Now())

	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, capa)
	assert.Equal(t, uint(5), replannedScale1)

	// WHEN
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale, maxScale, nil)
	scaleAdjustmentsPlanned.EXPECT().Set(float64(5))
	mocks.scaleAdjustments.EXPECT().WithLabelValues("planned").Return(scaleAdjustmentsPlanned)
	scaleAdjustmentsAdjusted.EXPECT().Set(float64(5))
	mocks.scaleAdjustments.EXPECT().WithLabelValues("adjusted").Return(scaleAdjustmentsAdjusted)

	replannedScale2 := capa.adjustPlanAccordingToSchedule(plannedScale, time.Now())

	// THEN
	assert.Equal(t, uint(5), replannedScale2)
}
