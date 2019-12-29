package capacityPlanner

import (
	"time"

	"github.com/thomasobenaus/sokar/helper"
)

func (cp *CapacityPlanner) adjustPlanAccordingToSchedule(currentlyPlannedScale uint, now time.Time) uint {
	const labelPlannedScale string = "planned"
	const labelAdjustedScale string = "adjusted"

	plannedScale := currentlyPlannedScale

	day := now.Weekday()
	hour := now.Hour()
	minute := now.Minute()
	at, err := helper.NewTime(uint(hour), uint(minute))

	if cp.schedule == nil {
		cp.logger.Debug().Msgf("No further adjustment of planned scale needed. No scaling schedule has been specified at [%s %s].", day, at)
		return plannedScale
	}

	// should never happen
	if err != nil {
		cp.logger.Warn().Msgf("Could not evaluate scaling schedule: %s.", err.Error())
		return plannedScale
	}

	minScale, maxScale, err := cp.schedule.ScaleRangeAt(day, at)
	if err != nil {
		cp.logger.Debug().Msgf("No further adjustment of planned scale needed. No scaling schedule entry found at current time [%s %s].", day, at)
		return plannedScale
	}

	plannedScale = fitIntoScaleRange(plannedScale, minScale, maxScale)
	if plannedScale != currentlyPlannedScale {
		cp.logger.Info().Msgf("The planned scale %d was adjusted to %d. Based on currently active schedule [%s %s] the scale has to be at least %d and at most %d.", currentlyPlannedScale, plannedScale, day, at, minScale, maxScale)
	} else {
		cp.logger.Debug().Msgf("No further adjustment of planned scale needed. The current scale of %d fits into the range of [%d-%d] of the currently active schedule [%s %s].", plannedScale, minScale, maxScale, day, at)
	}

	cp.metrics.scaleAdjustments.WithLabelValues(labelPlannedScale).Set(float64(currentlyPlannedScale))
	cp.metrics.scaleAdjustments.WithLabelValues(labelAdjustedScale).Set(float64(plannedScale))
	return plannedScale
}

func fitIntoScaleRange(scale, minScale, maxScale uint) uint {

	if scale < minScale {
		return minScale
	}

	if scale > maxScale {
		return maxScale
	}

	return scale
}
