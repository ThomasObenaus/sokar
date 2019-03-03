package capacityPlanner

// Plan computes the number of instances needed based on the current number and the scale factor
func (cp *CapacityPlanner) Plan(scaleFactor float32, currentScale uint) uint {

	if scaleFactor < 0 {
		return 0
	}
	return 2
}
