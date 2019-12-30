package sokar

import (
	"time"

	"github.com/thomasobenaus/sokar/helper"
)

// ScaleSchedule is an interface that is used to control when the ScaleScheduler should issue an alert
type ScaleSchedule interface {
	ScaleRangeAt(day time.Weekday, at helper.SimpleTime) (min uint, max uint, err error)
	IsActiveAt(day time.Weekday, at helper.SimpleTime) bool
}
