package capacityPlanner

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_capacityPlanner "github.com/thomasobenaus/sokar/test/capacityPlanner"
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
	scheduleIF := mock_capacityPlanner.NewMockScaleSchedule(mockCtrl)
	capa, err := New(Schedule(scheduleIF))
	minScale := uint(1)
	maxScale := uint(10)
	plannedScale1 := uint(0)
	plannedScale2 := uint(11)

	// WHEN
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale, maxScale, nil)
	replannedScale1 := capa.adjustPlanAccordingToSchedule(plannedScale1, time.Now())
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale, maxScale, nil)
	replannedScale2 := capa.adjustPlanAccordingToSchedule(plannedScale2, time.Now())

	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, capa)
	assert.Equal(t, uint(1), replannedScale1)
	assert.Equal(t, uint(10), replannedScale2)
}

func Test_ShouldNotAdjustScale(t *testing.T) {

	// GIVEN
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	scheduleIF := mock_capacityPlanner.NewMockScaleSchedule(mockCtrl)
	capa, err := New(Schedule(scheduleIF))
	plannedScale1 := uint(5)
	minScale1 := uint(1)
	maxScale1 := uint(10)
	err1 := fmt.Errorf("some err")

	// WHEN
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale1, maxScale1, err1)
	replannedScale1 := capa.adjustPlanAccordingToSchedule(plannedScale1, time.Now())
	scheduleIF.EXPECT().ScaleRangeAt(gomock.Any(), gomock.Any()).Return(minScale1, maxScale1, nil)
	replannedScale2 := capa.adjustPlanAccordingToSchedule(plannedScale1, time.Now())
	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, capa)
	assert.Equal(t, uint(5), replannedScale1)
	assert.Equal(t, uint(5), replannedScale2)
}
