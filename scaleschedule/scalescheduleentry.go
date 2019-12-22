package scaleschedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron"
)

type ScaleScheduleEntry struct {
	minScale uint
	maxScale uint

	// ScheduleDescription defines from which to which point in time at the
	// specified min- and max scale are demanded.
	// During this time range the ScaleScheduleEntry is relevant/ active. Outside of this time
	// range it is not relevant/ active.
	// The description has to be specified using cron notation (see: https://en.wikipedia.org/wiki/Cron).
	// Seconds, minutes, hours and day of week have to be specified.
	// e.g. "* * 8-9 MON-FRI" ==> 8am...10am at business days
	// Within the description a time range has to be specified. A singe reoccuring point in time is not valid.
	scheduleDescription string

	// this is the parsed variant and thus valid version of the scheduleDescription
	cron.Schedule
}

func NewScaleScheduleEntry(scheduleDescription string, minScale, maxScale uint) (ScaleScheduleEntry, error) {

	if maxScale <= minScale {
		return ScaleScheduleEntry{}, fmt.Errorf("MinScale (%d) has to be less then MaxScale (%d)", minScale, maxScale)
	}

	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dow)
	schedule, err := p.Parse(scheduleDescription)
	if err != nil {
		return ScaleScheduleEntry{}, err
	}

	// ensure that a time range is specified
	if !strings.Contains(scheduleDescription, "-") {
		return ScaleScheduleEntry{}, fmt.Errorf("The scheduleDescription (%s) has to be a time range and thus contain at least one '-'", scheduleDescription)
	}

	return ScaleScheduleEntry{
		minScale:            minScale,
		maxScale:            maxScale,
		scheduleDescription: scheduleDescription,
		Schedule:            schedule,
	}, nil
}

func calculateDuration(cronSpec string) (time.Duration, error) {

	// verify that the cron spec is valid
	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dow)
	_, err := p.Parse(cronSpec)
	if err != nil {
		return time.Second * 0, err
	}

	specParts := strings.Split(cronSpec, " ")
	if len(specParts) != 4 {
		return time.Second * 0, fmt.Errorf("Expected 4 parts in cronspec (%s) but got %d", cronSpec, len(cronSpec))
	}

	dowSpec := specParts[3]
	hourSpec := specParts[2]
	minSpec := specParts[1]
	secSpec := specParts[0]

}
