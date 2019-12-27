package scaleschedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	endTime1, _ := helper.NewTime(5, 0)
	startTime2, _ := helper.NewTime(23, 0)
	endTime2, _ := helper.NewTime(23, 30)
	startTime3, _ := helper.NewTime(7, 45)
	endTime3, _ := helper.NewTime(8, 15)
	minScale1 := uint(1)
	maxScale1 := uint(2)
	minScale2 := uint(11)
	maxScale2 := uint(22)
	minScale3 := uint(111)
	maxScale3 := uint(222)
	day := time.Monday
	scaleSchedule := New()

	// WHEN
	err1 := scaleSchedule.Insert(day, startTime1, endTime1, minScale1, maxScale1)
	err2 := scaleSchedule.Insert(day, startTime2, endTime2, minScale2, maxScale2)
	err3 := scaleSchedule.Insert(day, startTime3, endTime3, minScale3, maxScale3)

	//THEN
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	entries, ok := scaleSchedule.scheduleByDay[day]
	assert.NotNil(t, entries)
	assert.True(t, ok)
	assert.Len(t, entries, 3)
	assert.Equal(t, uint(0), entries[0].startMinute)
	assert.Equal(t, uint(300), entries[0].endMinute)
	assert.Equal(t, uint(1), entries[0].minScale)
	assert.Equal(t, uint(2), entries[0].maxScale)
	assert.Equal(t, uint(465), entries[1].startMinute)
	assert.Equal(t, uint(495), entries[1].endMinute)
	assert.Equal(t, uint(111), entries[1].minScale)
	assert.Equal(t, uint(222), entries[1].maxScale)
	assert.Equal(t, uint(1380), entries[2].startMinute)
	assert.Equal(t, uint(1410), entries[2].endMinute)
	assert.Equal(t, uint(11), entries[2].minScale)
	assert.Equal(t, uint(22), entries[2].maxScale)

	// GIVEN
	startTime, _ := helper.NewTime(0, 0)
	endTime, _ := helper.NewTime(0, 0)
	minScale := uint(4)
	maxScale := uint(44)
	day = time.Wednesday

	// WHEN
	err := scaleSchedule.Insert(day, startTime, endTime, minScale, maxScale)

	//THEN
	assert.NoError(t, err)
	entries, ok = scaleSchedule.scheduleByDay[day]
	assert.NotNil(t, entries)
	assert.True(t, ok)
	assert.Len(t, entries, 1)
	assert.Equal(t, uint(0), entries[0].startMinute)
	assert.Equal(t, uint(1440), entries[0].endMinute)
	assert.Equal(t, uint(4), entries[0].minScale)
	assert.Equal(t, uint(44), entries[0].maxScale)
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

func Test_InsertShouldFailOnConflict(t *testing.T) {

	// GIVEN
	startTime1, _ := helper.NewTime(0, 0)
	endTime1, _ := helper.NewTime(0, 0)
	day := time.Monday
	scaleSchedule := New()
	err := scaleSchedule.Insert(day, startTime1, endTime1, 1, 1)
	require.NoError(t, err)
	startTime2, _ := helper.NewTime(11, 0)
	endTime2, _ := helper.NewTime(12, 0)

	// WHEN
	err = scaleSchedule.Insert(day, startTime2, endTime2, 1, 1)

	//THEN
	assert.Error(t, err)
}

func Test_AtShouldFindSomething(t *testing.T) {

	// GIVEN
	scaleSchedule := New()
	day := time.Monday
	startTime, _ := helper.NewTime(0, 0)
	endTime, _ := helper.NewTime(0, 0)
	minScale := uint(1)
	maxScale := uint(2)
	err := scaleSchedule.Insert(day, startTime, endTime, minScale, maxScale)
	require.NoError(t, err)

	day = time.Wednesday
	startTime, _ = helper.NewTime(1, 0)
	endTime, _ = helper.NewTime(10, 0)
	minScale = uint(11)
	maxScale = uint(22)
	err = scaleSchedule.Insert(day, startTime, endTime, minScale, maxScale)
	require.NoError(t, err)

	// WHEN
	at1, _ := helper.NewTime(10, 15)
	entry1, err1 := scaleSchedule.at(time.Monday, at1)
	at2, _ := helper.NewTime(9, 15)
	entry2, err2 := scaleSchedule.at(time.Wednesday, at2)

	//THEN
	assert.NoError(t, err1)
	assert.Equal(t, uint(1), entry1.minScale)
	assert.NoError(t, err2)
	assert.Equal(t, uint(11), entry2.minScale)
}

func Test_AtShouldFindNothing(t *testing.T) {

	// GIVEN
	scaleSchedule := New()
	day := time.Monday
	startTime, _ := helper.NewTime(0, 0)
	endTime, _ := helper.NewTime(10, 0)
	minScale := uint(1)
	maxScale := uint(2)
	err := scaleSchedule.Insert(day, startTime, endTime, minScale, maxScale)
	require.NoError(t, err)

	// WHEN
	at, _ := helper.NewTime(10, 1)
	_, err1 := scaleSchedule.at(time.Monday, at)
	at, _ = helper.NewTime(9, 15)
	_, err2 := scaleSchedule.at(time.Wednesday, at)

	//THEN
	assert.Error(t, err1)
	assert.Error(t, err2)
}

func Test_ShouldReportConflictOnOverlap(t *testing.T) {

	// GIVEN
	scaleSchedule := New()
	day := time.Monday
	startTime, _ := helper.NewTime(9, 0)
	endTime, _ := helper.NewTime(10, 0)
	minScale := uint(1)
	maxScale := uint(2)
	err := scaleSchedule.Insert(day, startTime, endTime, minScale, maxScale)
	require.NoError(t, err)

	startTime, _ = helper.NewTime(9, 59)
	endTime, _ = helper.NewTime(10, 0)
	entry1 := entry{startTime.Minutes(), endTime.Minutes(), 0, 0}
	startTime, _ = helper.NewTime(8, 0)
	endTime, _ = helper.NewTime(9, 1)
	entry2 := entry{startTime.Minutes(), endTime.Minutes(), 0, 0}
	startTime, _ = helper.NewTime(9, 1)
	endTime, _ = helper.NewTime(9, 1)
	entry3 := entry{startTime.Minutes(), endTime.Minutes(), 0, 0}
	startTime, _ = helper.NewTime(0, 0)
	endTime, _ = helper.NewTime(23, 59)
	entry4 := entry{startTime.Minutes(), endTime.Minutes(), 0, 0}

	// WHEN
	hasConflict1 := scaleSchedule.isConflicting(time.Monday, entry1)
	hasConflict2 := scaleSchedule.isConflicting(time.Monday, entry2)
	hasConflict3 := scaleSchedule.isConflicting(time.Monday, entry3)
	hasConflict4 := scaleSchedule.isConflicting(time.Monday, entry4)

	//THEN
	assert.True(t, hasConflict1)
	assert.True(t, hasConflict2)
	assert.True(t, hasConflict3)
	assert.True(t, hasConflict4)
}

func Test_ShouldNotReportConflict(t *testing.T) {

	// GIVEN
	scaleSchedule := New()
	day := time.Monday
	startTime, _ := helper.NewTime(9, 0)
	endTime, _ := helper.NewTime(10, 0)
	minScale := uint(1)
	maxScale := uint(2)
	err := scaleSchedule.Insert(day, startTime, endTime, minScale, maxScale)
	require.NoError(t, err)

	startTime, _ = helper.NewTime(8, 58)
	endTime, _ = helper.NewTime(8, 59)
	entry1 := entry{startTime.Minutes(), endTime.Minutes(), 0, 0}
	startTime, _ = helper.NewTime(10, 1)
	endTime, _ = helper.NewTime(10, 2)
	entry2 := entry{startTime.Minutes(), endTime.Minutes(), 0, 0}
	startTime, _ = helper.NewTime(20, 0)
	endTime, _ = helper.NewTime(21, 0)
	entry3 := entry{startTime.Minutes(), endTime.Minutes(), 0, 0}

	// WHEN
	hasConflict1 := scaleSchedule.isConflicting(time.Monday, entry1)
	hasConflict2 := scaleSchedule.isConflicting(time.Monday, entry2)
	hasConflict3 := scaleSchedule.isConflicting(time.Monday, entry3)

	//THEN
	assert.False(t, hasConflict1)
	assert.False(t, hasConflict2)
	assert.False(t, hasConflict3)
}
