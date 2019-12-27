package scaleschedule

import (
	"fmt"
	"time"

	"github.com/thomasobenaus/sokar/helper"
)

type scheduleByDay map[time.Weekday][]*entry

// Schedule is a structure for creating and handling a scaling schedule
type Schedule struct {
	scheduleByDay scheduleByDay
}

type entry struct {
	startMinute uint
	endMinute   uint

	minScale uint
	maxScale uint
}

// New creates a new empty Schedule
func New() Schedule {
	return Schedule{
		scheduleByDay: make(scheduleByDay, 0),
	}
}

// Insert will insert a new scaling schedule entry into the Schedule. The entry will be formed based on the
// given parameters. Internally the parameters will be validated.
// In case the parameters are invalid or the entry to be inserted overlaps with an entry already present in the
// Schedule, an error will be returned.
func (s *Schedule) Insert(day time.Weekday, start, end helper.SimpleTime, minScale, maxScale uint) error {

	entries, ok := s.scheduleByDay[day]
	if !ok {
		s.scheduleByDay[day] = make([]*entry, 0)
		entries = s.scheduleByDay[day]
	}

	startMinutes := start.Minutes()
	endMinutes := end.Minutes()

	// catch the case 0:00 for endTime
	if endMinutes == 0 {
		endMinutes = 60 * 24
	}

	// ensure that start and end are valid
	if startMinutes >= endMinutes {
		return fmt.Errorf("StartTime (%s) has to be before endTime (%s)", start, end)
	}

	entry := entry{startMinutes, endMinutes, minScale, maxScale}
	s.scheduleByDay[day] = append(entries, &entry)

	// TODO: Sort and verify
	return nil
}
