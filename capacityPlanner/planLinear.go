package capacityPlanner

import (
	"math"

	"github.com/thomasobenaus/sokar/helper"
)

// planLinear increases/ decreases the given scale linearly based on the given scaleFactor.
// Therefore the scaleFactor is directly used to scale the number of currentScale by multiplication.
// newScale = (1+scaleFactor) * currentScale
func (cp *CapacityPlanner) planLinear(scaleFactor float32, currentScale uint) uint {

	if scaleFactor == 0 {
		return currentScale
	}

	increment := float64(scaleFactor * float32(currentScale) * cp.linearMode.ScaleFactorWeight)
	incrementInt := int(math.Ceil(increment))

	// at least scale up by one if scaleFactor is positive
	if incrementInt == 0 && scaleFactor > 0 {
		incrementInt = 1
	}

	// at least scale down by one if scaleFactor is negative
	if incrementInt == 0 && scaleFactor < 0 {
		incrementInt = -1
	}

	plannedScale := helper.IncUint(currentScale, incrementInt)

	return plannedScale
}
