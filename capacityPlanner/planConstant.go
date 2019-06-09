package capacityPlanner

import "github.com/thomasobenaus/sokar/helper"

// planConstant increases/ decreases the given scale by the given constant offset, regarding the given scaleFactor.
func (cp *CapacityPlanner) planConstant(scaleFactor float32, currentScale uint, offset uint) uint {
	plannedScale := currentScale
	if scaleFactor > 0 {
		plannedScale = currentScale + offset
	} else if scaleFactor < 0 && currentScale > 0 {
		plannedScale = helper.SubUint2(currentScale, offset)
	}
	return plannedScale
}
