package capacityPlanner

import "time"

// Plan computes the number of instances needed based on the current number and the scale factor
func (cp *CapacityPlanner) Plan(scaleFactor float32, currentScale uint) uint {

	plannedScale := uint(0)

	if cp.mode == CapaPlanningModeConstant {
		plannedScale = cp.planConstant(scaleFactor, currentScale, cp.offsetConstantMode)
	} else if cp.mode == CapaPlanningModeLinear {
		plannedScale = cp.planLinear(scaleFactor, currentScale)
	} else {
		cp.logger.Error().Msgf("Unknown planning mode %v. No planning done.", cp.mode)
	}

	cp.logger.Info().Msgf("Plan mode=%v, sf=%f, cs=%d, ps=%d.", cp.mode, scaleFactor, currentScale, plannedScale)
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
