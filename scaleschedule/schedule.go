package scaleschedule

import "time"

type ScaleSchedule interface {
	GetActiveScale(at time.Time) (active bool, max, min uint)
}

func New(scheduleEntries []ScaleScheduleEntry) (ScaleSchedule, error) {
	return nil, nil
}
