package capacityplanner

import (
	"time"

	"github.com/thomasobenaus/sokar/helper"
)

// ScaleSchedule is an interface that is used to control when the CapacityPlanner shall scale according to a ScaleSchedule
type ScaleSchedule interface {
	ScaleRangeAt(day time.Weekday, at helper.SimpleTime) (min uint, max uint, err error)
}
