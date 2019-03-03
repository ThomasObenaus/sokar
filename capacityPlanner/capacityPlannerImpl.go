package capacityPlanner

// Plan computes the number of instances needed based on the current number and the scale factor
func (cp *CapacityPlanner) Plan(scaleFactor float32, currentScale uint) uint {
	plannedScale := uint(0)
	// HACK: map scaleFactor directly to n - 1 or n + 1
	if scaleFactor > 0 {
		plannedScale = currentScale + 1
	} else if scaleFactor < 0 && currentScale > 0 {
		plannedScale = currentScale - 1
	}

	cp.logger.Info().Msgf("Plan sf=%f, cs=%d, ps=%d.", scaleFactor, currentScale, plannedScale)
	return plannedScale
}
