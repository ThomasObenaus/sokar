package scaleschedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/helper"
)

func Test_NewShouldInitiateEntriesEmpty(t *testing.T) {

	// GIVEN

	// WHEN
	scaleSchedule := New()

	//THEN
	assert.NotNil(t, scaleSchedule.scheduleByDay)
}

func Test_InsertShouldAddEntry(t *testing.T) {

	// GIVEN
	startTime1, _ := helper.NewTime(0, 0)
	endTime1, _ := helper.NewTime(0, 0)
	startTime2, _ := helper.NewTime(23, 0)
	endTime2, _ := helper.NewTime(23, 30)
	minScale1 := uint(1)
	maxScale1 := uint(2)
	minScale2 := uint(11)
	maxScale2 := uint(22)
	day := time.Monday
	scaleSchedule := New()

	// WHEN
	err1 := scaleSchedule.Insert(day, startTime1, endTime1, minScale1, maxScale1)
	err2 := scaleSchedule.Insert(day, startTime2, endTime2, minScale2, maxScale2)

	//THEN
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	entries, ok := scaleSchedule.scheduleByDay[day]
	assert.NotNil(t, entries)
	assert.True(t, ok)
	assert.Len(t, entries, 2)
	assert.Equal(t, uint(0), entries[0].startMinute)
	assert.Equal(t, uint(1440), entries[0].endMinute)
	assert.Equal(t, uint(1), entries[0].minScale)
	assert.Equal(t, uint(2), entries[0].maxScale)
	assert.Equal(t, uint(1380), entries[1].startMinute)
	assert.Equal(t, uint(1410), entries[1].endMinute)
	assert.Equal(t, uint(11), entries[1].minScale)
	assert.Equal(t, uint(22), entries[1].maxScale)

}

func Test_InsertShouldFailIfTimesAreNotValid(t *testing.T) {

	// GIVEN
	startTime, _ := helper.NewTime(1, 0)
	endTime, _ := helper.NewTime(0, 1)
	minScale := uint(1)
	maxScale := uint(2)
	day := time.Monday
	scaleSchedule := New()

	// WHEN
	err := scaleSchedule.Insert(day, startTime, endTime, minScale, maxScale)

	//THEN
	assert.Error(t, err)
}
