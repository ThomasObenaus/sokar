package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldParseScaleScheduleEntry(t *testing.T) {

	// GIVEN
	spec := "MON-FRI 2 4:45 10-30"

	// WHEN
	sEntry, err := parseScaleScheduleEntry(spec)

	//THEN
	assert.NoError(t, err)
	assert.Len(t, sEntry.Days, 5)
	assert.Equal(t, uint(2), sEntry.StartTime.Hour)
	assert.Equal(t, uint(0), sEntry.StartTime.Minute)
	assert.Equal(t, uint(4), sEntry.EndTime.Hour)
	assert.Equal(t, uint(45), sEntry.EndTime.Minute)
	assert.Equal(t, 10, sEntry.MinScale)
	assert.Equal(t, 30, sEntry.MaxScale)
}

func Test_ShouldNotParseScaleScheduleEntry(t *testing.T) {

	// GIVEN
	specEmpty := ""
	specParts := " 1 2 3 4 5"
	specInvalidDays := "invalid * * *"
	specInvalidStart := "* invalid * *"
	specInvalidEnd := "* * invalid *"
	specInvalidRange := "* * * invalid"

	// WHEN
	_, errEmpty := parseScaleScheduleEntry(specEmpty)
	_, errParts := parseScaleScheduleEntry(specParts)
	_, errInvalidDays := parseScaleScheduleEntry(specInvalidDays)
	_, errInvalidStart := parseScaleScheduleEntry(specInvalidStart)
	_, errInvalidEnd := parseScaleScheduleEntry(specInvalidEnd)
	_, errInvalidRange := parseScaleScheduleEntry(specInvalidRange)

	//THEN
	assert.Error(t, errEmpty)
	assert.Error(t, errParts)
	assert.Error(t, errInvalidDays)
	assert.Error(t, errInvalidStart)
	assert.Error(t, errInvalidEnd)
	assert.Error(t, errInvalidRange)
}

func Test_ShouldParseScaleRangeWildCard(t *testing.T) {
	// GIVEN
	scaleRangeSpecAllUnbound := "*"
	scaleRangeSpecMinUnbound := "*-10"
	scaleRangeSpecMaxUnbound := "2-*"

	// WHEN
	minAllU, maxAllU, errAllU := parseScaleRange(scaleRangeSpecAllUnbound)
	minMinU, maxMinU, errMinU := parseScaleRange(scaleRangeSpecMinUnbound)
	minMaxU, maxMaxU, errMaxU := parseScaleRange(scaleRangeSpecMaxUnbound)

	//THEN
	assert.NoError(t, errAllU)
	assert.NoError(t, errMinU)
	assert.NoError(t, errMaxU)
	assert.Equal(t, -1, minAllU)
	assert.Equal(t, -1, minMinU)
	assert.Equal(t, 2, minMaxU)
	assert.Equal(t, -1, maxAllU)
	assert.Equal(t, 10, maxMinU)
	assert.Equal(t, -1, maxMaxU)
}

func Test_ShouldParseScaleRange(t *testing.T) {
	// GIVEN
	scaleRangeSpec := "2-10"

	// WHEN
	min, max, err := parseScaleRange(scaleRangeSpec)

	//THEN
	assert.NoError(t, err)
	assert.Equal(t, 2, min)
	assert.Equal(t, 10, max)
}

func Test_ShouldParseNotScaleRange(t *testing.T) {
	// GIVEN
	scaleRangeSpecParts := "2-10-1"
	scaleRangeSpecMinNoNumber := "x-10"
	scaleRangeSpecMaxNoNumber := "2-x"
	scaleRangeSpecEmpty := ""
	scaleRangeSpecInvalidEntry := "s"
	scaleRangeSpecMaxGreaterThanMin := "22-1"
	scaleRangeSpecNegMin := "-1-5"
	scaleRangeSpecNegMax := "1--5"

	// WHEN
	_, _, errParts := parseScaleRange(scaleRangeSpecParts)
	_, _, errMinNoNumber := parseScaleRange(scaleRangeSpecMinNoNumber)
	_, _, errMaxNoNumber := parseScaleRange(scaleRangeSpecMaxNoNumber)
	_, _, errEmpty := parseScaleRange(scaleRangeSpecEmpty)
	_, _, errInvalidEntry := parseScaleRange(scaleRangeSpecInvalidEntry)
	_, _, errMaxGreaterThanMin := parseScaleRange(scaleRangeSpecMaxGreaterThanMin)
	_, _, errNegMin := parseScaleRange(scaleRangeSpecNegMin)
	_, _, errNegMax := parseScaleRange(scaleRangeSpecNegMax)

	//THEN
	assert.Error(t, errParts)
	assert.Error(t, errMinNoNumber)
	assert.Error(t, errMaxNoNumber)
	assert.Error(t, errEmpty)
	assert.Error(t, errInvalidEntry)
	assert.Error(t, errMaxGreaterThanMin)
	assert.Error(t, errNegMin)
	assert.Error(t, errNegMax)
}

func Test_ShouldParseTimeSpecWildCard(t *testing.T) {
	// GIVEN
	timeSpec := "*"

	// WHEN
	parsedTime, err := parseTime(timeSpec)

	//THEN
	assert.NoError(t, err)
	assert.Equal(t, uint(0), parsedTime.Hour)
	assert.Equal(t, uint(0), parsedTime.Minute)
}

func Test_ShouldParseTimeSpecHoursOnly(t *testing.T) {
	// GIVEN
	timeSpec0 := "0"
	timeSpec1 := "23"

	// WHEN
	parsedTime0, err0 := parseTime(timeSpec0)
	parsedTime1, err1 := parseTime(timeSpec1)

	//THEN
	assert.NoError(t, err0)
	assert.NoError(t, err1)
	assert.Equal(t, uint(0), parsedTime0.Hour)
	assert.Equal(t, uint(23), parsedTime1.Hour)
	assert.Equal(t, uint(0), parsedTime0.Minute)
	assert.Equal(t, uint(0), parsedTime1.Minute)
}

func Test_ShouldParseTimeSpec(t *testing.T) {
	// GIVEN
	timeSpec := "13:15"

	// WHEN
	parsedTime, err := parseTime(timeSpec)

	//THEN
	assert.NoError(t, err)
	assert.Equal(t, uint(13), parsedTime.Hour)
	assert.Equal(t, uint(15), parsedTime.Minute)
}

func Test_ShouldNotParseTimeSpecOutOfRange(t *testing.T) {
	// GIVEN
	timeSpecHourNeg := "-1:15"
	timeSpecHourOnlyNeg := "-1"
	timeSpecMinNeg := "1:-15"
	timeSpecHourOoR := "24:15"
	timeSpecHourOnlyOoR := "24"
	timeSpecMinOoR := "23:60"

	// WHEN
	_, errHourNeg := parseTime(timeSpecHourNeg)
	_, errHourOnlyNeg := parseTime(timeSpecHourOnlyNeg)
	_, errMinNeg := parseTime(timeSpecMinNeg)
	_, errHourOoR := parseTime(timeSpecHourOoR)
	_, errHourOnlyOoR := parseTime(timeSpecHourOnlyOoR)
	_, errMinOoR := parseTime(timeSpecMinOoR)

	//THEN
	assert.Error(t, errHourNeg)
	assert.Error(t, errHourOnlyNeg)
	assert.Error(t, errMinNeg)
	assert.Error(t, errHourOoR)
	assert.Error(t, errHourOnlyOoR)
	assert.Error(t, errMinOoR)
}

func Test_ShouldNotParseEmptyTimeSpec(t *testing.T) {
	// GIVEN
	timeSpec := ""

	// WHEN
	_, err := parseTime(timeSpec)

	//THEN
	assert.Error(t, err)
}

func Test_ShouldNotParseMalformedTimeSpec(t *testing.T) {
	// GIVEN
	timeSpec := "12:12:12"

	// WHEN
	_, err := parseTime(timeSpec)

	//THEN
	assert.Error(t, err)
}

func Test_ShouldParseEmptyDaysSpec(t *testing.T) {
	// GIVEN
	daysSpec := ""

	// WHEN
	days, err := parseDays(daysSpec)

	//THEN
	assert.NoError(t, err)
	assert.Empty(t, days)
}

func Test_ShouldParseWildcard(t *testing.T) {
	// GIVEN
	daysSpec := "*"

	// WHEN
	days, err := parseDays(daysSpec)

	//THEN
	assert.NoError(t, err)
	assert.Len(t, days, 7)
	assert.Equal(t, time.Sunday, days[0])
	assert.Equal(t, time.Monday, days[1])
	assert.Equal(t, time.Tuesday, days[2])
	assert.Equal(t, time.Wednesday, days[3])
	assert.Equal(t, time.Thursday, days[4])
	assert.Equal(t, time.Friday, days[5])
	assert.Equal(t, time.Saturday, days[6])
}

func Test_ShouldNotParseInvalidDay(t *testing.T) {
	// GIVEN
	daysSpec := "invalid"

	// WHEN
	days, err := parseDays(daysSpec)

	//THEN
	assert.Error(t, err)
	assert.Empty(t, days)
}

func Test_ShouldParseSingleDaysSpec(t *testing.T) {
	// GIVEN
	daysSpecMON := "MON"
	daysSpecTUE := "TUE"
	daysSpecWED := "WED"
	daysSpecTHU := "THU"
	daysSpecFRI := "FRI"
	daysSpecSAT := "SAT"
	daysSpecSUN := "SUN"

	// WHEN
	daysMON, errMON := parseDays(daysSpecMON)
	daysTUE, errTUE := parseDays(daysSpecTUE)
	daysWED, errWED := parseDays(daysSpecWED)
	daysTHU, errTHU := parseDays(daysSpecTHU)
	daysFRI, errFRI := parseDays(daysSpecFRI)
	daysSAT, errSAT := parseDays(daysSpecSAT)
	daysSUN, errSUN := parseDays(daysSpecSUN)

	//THEN
	assert.NoError(t, errMON)
	assert.NoError(t, errTUE)
	assert.NoError(t, errWED)
	assert.NoError(t, errTHU)
	assert.NoError(t, errFRI)
	assert.NoError(t, errSAT)
	assert.NoError(t, errSUN)
	assert.Len(t, daysMON, 1)
	assert.Len(t, daysTUE, 1)
	assert.Len(t, daysWED, 1)
	assert.Len(t, daysTHU, 1)
	assert.Len(t, daysFRI, 1)
	assert.Len(t, daysSAT, 1)
	assert.Len(t, daysSUN, 1)
	assert.Equal(t, time.Monday, daysMON[0])
	assert.Equal(t, time.Tuesday, daysTUE[0])
	assert.Equal(t, time.Wednesday, daysWED[0])
	assert.Equal(t, time.Thursday, daysTHU[0])
	assert.Equal(t, time.Friday, daysFRI[0])
	assert.Equal(t, time.Saturday, daysSAT[0])
	assert.Equal(t, time.Sunday, daysSUN[0])

	// GIVEN
	daysSpecMON = "1"
	daysSpecTUE = "2"
	daysSpecWED = "3"
	daysSpecTHU = "4"
	daysSpecFRI = "5"
	daysSpecSAT = "6"
	daysSpecSUN = "0"

	// WHEN
	daysMON, errMON = parseDays(daysSpecMON)
	daysTUE, errTUE = parseDays(daysSpecTUE)
	daysWED, errWED = parseDays(daysSpecWED)
	daysTHU, errTHU = parseDays(daysSpecTHU)
	daysFRI, errFRI = parseDays(daysSpecFRI)
	daysSAT, errSAT = parseDays(daysSpecSAT)
	daysSUN, errSUN = parseDays(daysSpecSUN)

	//THEN
	assert.NoError(t, errMON)
	assert.NoError(t, errTUE)
	assert.NoError(t, errWED)
	assert.NoError(t, errTHU)
	assert.NoError(t, errFRI)
	assert.NoError(t, errSAT)
	assert.NoError(t, errSUN)
	assert.Len(t, daysMON, 1)
	assert.Len(t, daysTUE, 1)
	assert.Len(t, daysWED, 1)
	assert.Len(t, daysTHU, 1)
	assert.Len(t, daysFRI, 1)
	assert.Len(t, daysSAT, 1)
	assert.Len(t, daysSUN, 1)
	assert.Equal(t, time.Monday, daysMON[0])
	assert.Equal(t, time.Tuesday, daysTUE[0])
	assert.Equal(t, time.Wednesday, daysWED[0])
	assert.Equal(t, time.Thursday, daysTHU[0])
	assert.Equal(t, time.Friday, daysFRI[0])
	assert.Equal(t, time.Saturday, daysSAT[0])
	assert.Equal(t, time.Sunday, daysSUN[0])
}

func Test_ShouldNotParseRangeDaysSpec(t *testing.T) {
	// GIVEN
	daysSpec := "INVALID-FRI"

	// WHEN
	days, err := parseDays(daysSpec)

	//THEN
	assert.Error(t, err)
	assert.Empty(t, days)

	// GIVEN
	daysSpec = "-FRI"

	// WHEN
	days, err = parseDays(daysSpec)

	//THEN
	assert.Error(t, err)
	assert.Empty(t, days)

	// GIVEN
	daysSpec = "FRI-INVALID"

	// WHEN
	days, err = parseDays(daysSpec)

	//THEN
	assert.Error(t, err)
	assert.Empty(t, days)

	// GIVEN
	daysSpec = "FRI-"

	// WHEN
	days, err = parseDays(daysSpec)

	//THEN
	assert.Error(t, err)
	assert.Empty(t, days)
}

func Test_ShouldParseRangeDaysSpec(t *testing.T) {
	// GIVEN
	daysSpec := "MON-FRI"

	// WHEN
	days, err := parseDays(daysSpec)

	//THEN
	assert.NoError(t, err)
	assert.Len(t, days, 5)
	assert.Contains(t, days, time.Monday)
	assert.Contains(t, days, time.Tuesday)
	assert.Contains(t, days, time.Wednesday)
	assert.Contains(t, days, time.Thursday)
	assert.Contains(t, days, time.Friday)
	assert.NotContains(t, days, time.Saturday)
	assert.NotContains(t, days, time.Sunday)

	// GIVEN
	daysSpec = "FRI-MON"

	// WHEN
	days, err = parseDays(daysSpec)

	//THEN
	assert.NoError(t, err)
	assert.Len(t, days, 4)
	assert.Contains(t, days, time.Friday)
	assert.Contains(t, days, time.Saturday)
	assert.Contains(t, days, time.Sunday)
	assert.Contains(t, days, time.Monday)
	assert.NotContains(t, days, time.Tuesday)
	assert.NotContains(t, days, time.Wednesday)
	assert.NotContains(t, days, time.Thursday)

	// GIVEN
	daysSpec = "FRI-FRI"

	// WHEN
	days, err = parseDays(daysSpec)

	//THEN
	assert.NoError(t, err)
	assert.Len(t, days, 1)
	assert.Contains(t, days, time.Friday)
	assert.NotContains(t, days, time.Monday)
	assert.NotContains(t, days, time.Tuesday)
	assert.NotContains(t, days, time.Wednesday)
	assert.NotContains(t, days, time.Thursday)
	assert.NotContains(t, days, time.Saturday)
	assert.NotContains(t, days, time.Sunday)
}

//func Test_NewShouldSucceed(t *testing.T) {
//	// GIVEN
//	scheduleDescription := "* * 8-9 MON-FRI"
//	minScale := uint(1)
//	maxScale := uint(10)
//
//	// WHEN
//	sse, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)
//
//	//THEN
//	assert.NoError(t, err)
//	assert.Equal(t, minScale, sse.minScale)
//	assert.Equal(t, maxScale, sse.maxScale)
//	assert.Equal(t, scheduleDescription, sse.scheduleDescription)
//}
//
//func Test_NewShouldFailOnMalformedCron(t *testing.T) {
//	// GIVEN
//	scheduleDescription := "invalid"
//	minScale := uint(1)
//	maxScale := uint(10)
//
//	// WHEN
//	_, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)
//
//	//THEN
//	assert.Error(t, err)
//
//	// GIVEN
//	scheduleDescription = "* * 8-9"
//	minScale = uint(1)
//	maxScale = uint(10)
//
//	// WHEN
//	_, err = NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)
//
//	//THEN
//	assert.Error(t, err)
//}
//
//func Test_NewShouldFailOnInvalidScaleValues(t *testing.T) {
//	// GIVEN
//	scheduleDescription := "* * 8-9 MON-FRI"
//	minScale := uint(1)
//	maxScale := uint(1)
//
//	// WHEN
//	_, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)
//
//	//THEN
//	assert.Error(t, err)
//
//	// GIVEN
//	minScale = uint(2)
//	maxScale = uint(1)
//
//	// WHEN
//	_, err = NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)
//
//	//THEN
//	assert.Error(t, err)
//}
//
//func Test_NewShouldFailOnNoTimeRangeSpecified(t *testing.T) {
//	// GIVEN
//	scheduleDescription := "* * 8 FRI"
//	minScale := uint(1)
//	maxScale := uint(2)
//
//	// WHEN
//	_, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)
//
//	//THEN
//	assert.Error(t, err)
//}
//
//func Test_NewShouldFailOn(t *testing.T) {
//	scheduleDescription := "* ? * MON-FRI"
//	minScale := uint(1)
//	maxScale := uint(2)
//
//	// WHEN
//	sse1, err := NewScaleScheduleEntry(scheduleDescription, minScale, maxScale)
//	require.NoError(t, err)
//	n1 := sse1.Next(time.Now())
//	log.Printf("=> %v\n", n1)
//	n2 := sse1.Next(n1)
//	log.Printf("=> %v\n", n2)
//
//}
//
