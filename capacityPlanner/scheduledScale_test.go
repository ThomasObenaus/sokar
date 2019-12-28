package capacityPlanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
