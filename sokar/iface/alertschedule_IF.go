package sokar

import (
	"time"

	"github.com/thomasobenaus/sokar/helper"
)

// AlertSchedule is an interface that is used to control when the AlertScheduler should issue an alert
type AlertSchedule interface {
	IsActiveAt(day time.Weekday, at helper.SimpleTime) bool
}
