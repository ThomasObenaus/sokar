package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/thomasobenaus/sokar/helper"
)

// ScaleScheduleEntry represents one entry of a ScaleSchedule
type ScaleScheduleEntry struct {
	Days      []time.Weekday    `json:"days,omitempty"`
	StartTime helper.SimpleTime `json:"start_time,omitempty"`
	EndTime   helper.SimpleTime `json:"end_time,omitempty"`
	// MinScale -1 means unbound
	MinScale int `json:"min_scale,omitempty"`
	// MaxScale -1 means unbound
	MaxScale int `json:"max_scale,omitempty"`

	spec string
}

func (s ScaleScheduleEntry) String() string {
	return s.spec
}

func parseScaleScheduleEntry(spec string) (ScaleScheduleEntry, error) {
	spec = strings.TrimSpace(spec)
	if len(spec) == 0 {
		return ScaleScheduleEntry{}, fmt.Errorf("ScaleScheduleSpec (%s) is empty", spec)
	}

	parts := strings.Split(spec, " ")
	if len(parts) != 4 {
		return ScaleScheduleEntry{}, fmt.Errorf("ScaleScheduleSpec (%s) is malformed", spec)
	}
	daysSpec := parts[0]
	startTimeSpec := parts[1]
	endTimeSpec := parts[2]
	scaleRangeSpec := parts[3]

	days, err := parseDays(daysSpec)
	if err != nil {
		return ScaleScheduleEntry{}, err
	}

	startTime, err := parseTime(startTimeSpec)
	if err != nil {
		return ScaleScheduleEntry{}, err
	}

	endTime, err := parseTime(endTimeSpec)
	if err != nil {
		return ScaleScheduleEntry{}, err
	}

	min, max, err := parseScaleRange(scaleRangeSpec)
	if err != nil {
		return ScaleScheduleEntry{}, err
	}

	return ScaleScheduleEntry{Days: days, StartTime: startTime, EndTime: endTime, MinScale: min, MaxScale: max, spec: spec}, nil
}

func parseScaleRange(scaleRangeSpec string) (min int, max int, err error) {

	scaleRangeSpec = strings.TrimSpace(scaleRangeSpec)
	if len(scaleRangeSpec) == 0 {
		return 0, 0, fmt.Errorf("Scalespec (%s) is empty", scaleRangeSpec)
	}

	parts := strings.Split(scaleRangeSpec, "-")
	minStr := ""
	maxStr := ""

	switch len(parts) {
	// single element --> wildcard
	case 1:
		// wildcard --> -1,-1
		if strings.Compare(parts[0], "*") == 0 {
			return -1, -1, nil
		}
		return 0, 0, fmt.Errorf("Scalespec (%s) is malformed, min and max have to be specified", scaleRangeSpec)
	// multiple elements --> hours and minutes
	case 2:
		minStr = strings.TrimSpace(parts[0])
		maxStr = strings.TrimSpace(parts[1])
		break
	default:
		return 0, 0, fmt.Errorf("Scalespec (%s) is malformed", scaleRangeSpec)
	}

	// parse min
	minVal := int64(-1)
	if strings.Compare(minStr, "*") != 0 {
		minVal, err = strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("Scalespec (%s) is malformed. %s is not a number", scaleRangeSpec, minStr)
		}
		if minVal < 0 {
			return 0, 0, fmt.Errorf("Scalespec (%s) is malformed. Negative min values are not allowed", scaleRangeSpec)
		}
	}

	// parse max
	maxVal := int64(-1)
	if strings.Compare(maxStr, "*") != 0 {
		maxVal, err = strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("Scalespec (%s) is malformed. %s is not a number", scaleRangeSpec, maxStr)
		}
		if maxVal < 0 {
			return 0, 0, fmt.Errorf("Scalespec (%s) is malformed. Negative max values are not allowed", scaleRangeSpec)
		}
	}

	if maxVal < minVal && maxVal != -1 {
		return 0, 0, fmt.Errorf("Max scale value (%d) must not be less than the min (%d) scale value", maxVal, minVal)
	}

	return int(minVal), int(maxVal), nil
}

func parseTime(timeSpec string) (helper.SimpleTime, error) {

	timeSpec = strings.TrimSpace(timeSpec)
	if len(timeSpec) == 0 {
		return helper.SimpleTime{}, fmt.Errorf("Timespec (%s) is empty", timeSpec)
	}

	parts := strings.Split(timeSpec, ":")
	hourStr := ""
	minuteStr := ""

	switch len(parts) {
	// single element --> hours or wildcard
	case 1:
		// wildcard --> 0:00
		if strings.Compare(parts[0], "*") == 0 {
			return helper.SimpleTime{}, nil
		}
		hourStr = strings.TrimSpace(parts[0])
		break
	// multiple elements --> hours and minutes
	case 2:
		hourStr = strings.TrimSpace(parts[0])
		minuteStr = strings.TrimSpace(parts[1])
		break
	default:
		return helper.SimpleTime{}, fmt.Errorf("Timespec is malformed (%s)", timeSpec)
	}

	// parse hour
	hour, err := strconv.ParseUint(hourStr, 10, 64)
	if err != nil {
		return helper.SimpleTime{}, fmt.Errorf("Timespec is malformed. Hour (%s) is unknown", hourStr)
	}
	if hour < 0 || hour > 23 {
		return helper.SimpleTime{}, fmt.Errorf("Timespec is malformed. Hour (%s) is not between 0 and 23", hourStr)
	}

	// parse minute
	minute := uint64(0)

	if len(minuteStr) > 0 {
		minute, err = strconv.ParseUint(minuteStr, 10, 64)
		if err != nil {
			return helper.SimpleTime{}, fmt.Errorf("Timespec is malformed. Minute (%s) is unknown", minuteStr)
		}
		if minute < 0 || minute > 59 {
			return helper.SimpleTime{}, fmt.Errorf("Timespec is malformed. Minute (%s) is not between 0 and 59", minuteStr)
		}
	}

	return helper.NewTime(uint(hour), uint(minute))
}

func parseDays(daysSpec string) ([]time.Weekday, error) {
	dowTokens := map[string]time.Weekday{
		`0`: time.Sunday, `sun`: time.Sunday, `sunday`: time.Sunday,
		`1`: time.Monday, `mon`: time.Monday, `monday`: time.Monday,
		`2`: time.Tuesday, `tue`: time.Tuesday, `tuesday`: time.Tuesday,
		`3`: time.Wednesday, `wed`: time.Wednesday, `wednesday`: time.Wednesday,
		`4`: time.Thursday, `thu`: time.Thursday, `thursday`: time.Thursday,
		`5`: time.Friday, `fri`: time.Friday, `friday`: time.Friday,
		`6`: time.Saturday, `sat`: time.Saturday, `saturday`: time.Saturday,
	}

	days := make([]time.Weekday, 0)

	daysSpec = strings.TrimSpace(daysSpec)
	if len(daysSpec) == 0 {
		return make([]time.Weekday, 0), nil
	}

	parts := strings.Split(daysSpec, "-")

	// single day
	if len(parts) == 1 {
		dowKey := strings.ToLower(parts[0])

		// wildcard --> all weekdays
		if strings.Compare(dowKey, "*") == 0 {
			days = append(days, time.Sunday)
			days = append(days, time.Monday)
			days = append(days, time.Tuesday)
			days = append(days, time.Wednesday)
			days = append(days, time.Thursday)
			days = append(days, time.Friday)
			days = append(days, time.Saturday)
			return days, nil
		}

		dow, ok := dowTokens[dowKey]
		if !ok {
			return nil, fmt.Errorf("Unknown day '%s'", dowKey)
		}
		days = append(days, dow)
		return days, nil
	}

	// multiple days
	startDay := strings.ToLower(strings.TrimSpace(parts[0]))
	if len(startDay) == 0 {
		return nil, fmt.Errorf("Start day of daySpec range (%s) is empty", daysSpec)
	}

	endDay := strings.ToLower(strings.TrimSpace(parts[1]))
	if len(endDay) == 0 {
		return nil, fmt.Errorf("End day of daySpec range (%s) is empty", daysSpec)
	}

	dowStart, ok := dowTokens[startDay]
	if !ok {
		return nil, fmt.Errorf("Unknown start day '%s'", startDay)
	}
	dowEnd, ok := dowTokens[endDay]
	if !ok {
		return nil, fmt.Errorf("Unknown end day '%s'", endDay)
	}

	if dowStart < dowEnd {
		for i := dowStart; i <= dowEnd; i++ {
			dow, _ := dowTokens[fmt.Sprintf("%d", i)]
			days = append(days, dow)
		}
		return days, nil
	}

	// again, one single day
	if dowStart == dowEnd {
		days = append(days, dowStart)
		return days, nil
	}

	// dowStart >= dowEnd
	numDays := (6 - int(dowStart)) + int(dowEnd) + 2
	for i := 0; i < numDays; i++ {
		dow, _ := dowTokens[fmt.Sprintf("%d", (i+int(dowStart))%7)]
		days = append(days, dow)
	}

	return days, nil
}

// NewScaleScheduleEntry creates a new Entry based on the given specification
// The spec of an entry consist of five parts '<days> <start-time> <end-time> <scale-range>'.
// The parts are separated by '<space>'
// 1. <days>
//   - Specifies at which weekdays the schedule shall be active.
//   - Valid values: 'MON,TUE,WED,THU,FRI,SAT,SUN' and '0,1,2,3,4,5,6'
//   - Ranges are allowed e.g. 'MON-FRI', but have to be ascending (e.g. 'FRI-MON' is invalid)
//   - The wildcard '*' is also allowed and means any day.
// 2. <start-time>
//   - Specifies the time at which the schedule begins.
//   - Format: '<hour>:<minute>'. Where hour is a number between 0 and 23 and minute a number between 0 and 59.
//   - The minute qualifier is optional. Thus instead of specifying '13:00', '13' is sufficient.
//   - The wildcard '*' is also allowed and means start of the day (0:00 - midnight).
// 3. <end-time>
//   - Specifies the time at which the schedule ends.
//   - Format: '<hour>:<minute>'. Where hour is a number between 0 and 23 and minute a number between 0 and 59.
//   - The minute qualifier is optional. Thus instead of specifying '13:00', '13' is sufficient.
//   - The wildcard '*' is also allowed and means end of the day (0:00 - midnight).
// 4. <scale-range>
//   - Specifies the range within which the scale of the scale object shall be kept during this schedule.
//   - Format: '<min-scale>-<max-scale>', where both scale values are 'uint'.
//   - The wildcard '*' is also allowed and means unbound. For example '*-10' means the min scale is not bound whereas the max is set to 10.
//   - It is allowed to just specify '*' instead of '*-*' if both, min- and max-scale shall be unbound.
//     Even though it makes no sense, since no scheduled scaling would be done in this case.
func NewScaleScheduleEntry(spec string) (ScaleScheduleEntry, error) {
	return parseScaleScheduleEntry(spec)
}
