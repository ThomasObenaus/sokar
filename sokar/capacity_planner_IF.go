package sokar

// CapacityPlanner represents the interface for a component deciding
// the amount of instances needed of a job at a certain point in time.
type CapacityPlanner interface {
	// Plan plans how many instances are needed based on the given
	// scaleFactor
	Plan(scaleFactor float32, currentScale uint) uint
}
