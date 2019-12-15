package capacityPlanner

import "time"

// Plan computes the number of instances needed based on the current number and the scale factor
func (cp *CapacityPlanner) Plan(scaleFactor float32, currentScale uint) uint {

	plannedScale := uint(0)

	if cp.constantMode != nil {
		offset := cp.constantMode.Offset
		plannedScale = cp.planConstant(scaleFactor, currentScale, offset)
		cp.logger.Info().Msgf("Plan mode=constant, sf=%f, cs=%d, ps=%d, off=%d.", scaleFactor, currentScale, plannedScale, offset)
	} else if cp.linearMode != nil {
		plannedScale = cp.planLinear(scaleFactor, currentScale)
		cp.logger.Info().Msgf("Plan mode=linear, sf=%f, cs=%d, ps=%d sfW=%f.", scaleFactor, currentScale, plannedScale, cp.linearMode.ScaleFactorWeight)
	} else {
		cp.logger.Error().Msgf("No planning mode defined")
	}

	// TODO: Ensure here that the scheduled scale is satisfied.
	// TODO: Check how we know here that scheduled scaling is active?
	// TODO: Introduce a new event-source to ensure that (in scheduled scaling) the CAP can do the stuff described above
	// TODO: New Metrics?
	// TODO: Ensure that the old metrics still make sense (since we introduce a new event source)

	return plannedScale
}

// IsCoolingDown returns true if the CapacityPlanner thinks that currently a new scaling
// would not be a good idea.
func (cp *CapacityPlanner) IsCoolingDown(timeOfLastScale time.Time, scaleDown bool) (cooldownActive bool, cooldownTimeLeft time.Duration) {
	now := time.Now()

	dur := cp.upScaleCooldownPeriod
	if scaleDown {
		dur = cp.downScaleCooldownPeriod
	}

	// still cooling down
	if timeOfLastScale.Add(dur).After(now) {
		return true, timeOfLastScale.Add(dur).Sub(now)
	}

	// not cooling down any more
	return false, time.Second * 0
}
