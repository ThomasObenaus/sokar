package capacityPlanner

import "time"

// Plan computes the number of instances needed based on the current number and the scale factor
func (cp *CapacityPlanner) Plan(scaleFactor float32, currentScale uint) uint {

	plannedScale := uint(0)
	planMode := "undefined"

	if cp.constantMode != nil {
		planMode = "constant"
		plannedScale = cp.planConstant(scaleFactor, currentScale, cp.constantMode.Offset)
	} else if cp.linearMode != nil {
		planMode = "linear"
		plannedScale = cp.planLinear(scaleFactor, currentScale)
	} else {
		cp.logger.Error().Msgf("No planning mode defined")
	}

	cp.logger.Info().Msgf("Plan mode=%v, sf=%f, cs=%d, ps=%d.", planMode, scaleFactor, currentScale, plannedScale)
	return plannedScale
}

// IsCoolingDown returns true if the CapacityPlanner thinks that currently a new scaling
// would not be a good idea.
func (cp *CapacityPlanner) IsCoolingDown(timeOfLastScale time.Time, scaleDown bool) bool {
	now := time.Now()

	dur := cp.upScaleCooldownPeriod
	if scaleDown {
		dur = cp.downScaleCooldownPeriod
	}

	if timeOfLastScale.Add(dur).After(now) {
		return true
	}

	return false
}
