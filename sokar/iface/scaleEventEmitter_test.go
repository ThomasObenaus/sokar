package sokar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldCreateNewScaleEvent(t *testing.T) {

	// GIVEN
	// WHEN
	se1 := NewScaleEvent(1.2)
	se2 := NewScheduledScaleEvent()

	// THEN
	assert.Equal(t, float32(1.2), se1.ScaleFactor())
	assert.Equal(t, scaleEventRegular, se1.sType)
	assert.Equal(t, float32(0), se2.ScaleFactor())
	assert.Equal(t, scaleEventScheduled, se2.sType)
}
