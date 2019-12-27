package scaleschedule

import (
	"fmt"
	"sort"
	"time"

	"github.com/thomasobenaus/sokar/helper"
)

type scheduleByDay map[time.Weekday][]*entry

// Schedule is a structure for creating and handling a scaling schedule
type Schedule struct {
	scheduleByDay scheduleByDay
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
	if s.isConflicting(day, entry) {
		return fmt.Errorf("%s can't be inserted since it conflicts/ overlaps with entries of the schedule", entry)
	}

	s.scheduleByDay[day] = append(entries, &entry)

	// Sort the entries based on their start on that day
	sort.Sort(byStartMinute(s.scheduleByDay[day]))
	return nil
}

// isConflicting returns true in case the given entry overlaps with any entry of the current
// schedule at the specified day. If no conflict/ overlap is detected false will be returned.
func (s *Schedule) isConflicting(day time.Weekday, e entry) bool {

	entries, ok := s.scheduleByDay[day]
	if !ok {
		return false
	}

	for _, currentEntry := range entries {
		// should not happen
		if currentEntry == nil {
			continue
		}

		cStart := currentEntry.startMinute
		cEnd := currentEntry.endMinute

		if cStart >= e.startMinute && cEnd <= e.endMinute {
			return true
		}

		if cStart >= e.startMinute && cEnd >= e.endMinute && cStart <= e.endMinute {
			return true
		}

		if cStart <= e.startMinute && cEnd <= e.endMinute && cEnd >= e.startMinute {
			return true
		}

		if cStart <= e.startMinute && cEnd >= e.endMinute {
			return true
		}
	}

	return false
}

// Returns the scale schedule entry for the given day whose time range covers the given time.
// In case no entry can be found an error is returned.
func (s *Schedule) at(day time.Weekday, at helper.SimpleTime) (entry, error) {

	entries, ok := s.scheduleByDay[day]
	if !ok {
		return entry{}, fmt.Errorf("No entry at this day (%s)", day)
	}

	relevantMinute := at.Minutes()
	for _, currentEntry := range entries {

		// should not happen
		if currentEntry == nil {
			return entry{}, fmt.Errorf("The given entry is nil")
		}

		if currentEntry.startMinute <= relevantMinute && currentEntry.endMinute >= relevantMinute {
			return *currentEntry, nil
		}
	}
	return entry{}, fmt.Errorf("No entry found at %s %s", day, at)
}
