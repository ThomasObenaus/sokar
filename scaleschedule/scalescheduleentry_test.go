package scaleschedule

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewShouldSucceed(t *testing.T) {
	// GIVEN
	scheduleDescription := "* * 8-9 MON-FRI"
	minScale := uint(1)
	maxScale := uint(10)

	// WHEN
	sse, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)

	//THEN
	assert.NoError(t, err)
	assert.Equal(t, minScale, sse.minScale)
	assert.Equal(t, maxScale, sse.maxScale)
	assert.Equal(t, scheduleDescription, sse.scheduleDescription)
}

func Test_NewShouldFailOnMalformedCron(t *testing.T) {
	// GIVEN
	scheduleDescription := "invalid"
	minScale := uint(1)
	maxScale := uint(10)

	// WHEN
	_, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)

	//THEN
	assert.Error(t, err)

	// GIVEN
	scheduleDescription = "* * 8-9"
	minScale = uint(1)
	maxScale = uint(10)

	// WHEN
	_, err = NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)

	//THEN
	assert.Error(t, err)
}

func Test_NewShouldFailOnInvalidScaleValues(t *testing.T) {
	// GIVEN
	scheduleDescription := "* * 8-9 MON-FRI"
	minScale := uint(1)
	maxScale := uint(1)

	// WHEN
	_, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)

	//THEN
	assert.Error(t, err)

	// GIVEN
	minScale = uint(2)
	maxScale = uint(1)

	// WHEN
	_, err = NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)

	//THEN
	assert.Error(t, err)
}

func Test_NewShouldFailOnNoTimeRangeSpecified(t *testing.T) {
	// GIVEN
	scheduleDescription := "* * 8 FRI"
	minScale := uint(1)
	maxScale := uint(2)

	// WHEN
	_, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)

	//THEN
	assert.Error(t, err)
}

func Test_NewShouldFailOn(t *testing.T) {
	scheduleDescription := "* ? * MON-FRI"
	minScale := uint(1)
	maxScale := uint(2)

	// WHEN
	sse1, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)
	require.NoError(t, err)
	n1 := sse1.Next(time.Now())
	log.Printf("=> %v\n", n1)
	n2 := sse1.Next(n1)
	log.Printf("=> %v\n", n2)

}
