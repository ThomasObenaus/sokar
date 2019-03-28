package sokar

import "time"

// CapacityPlanner represents the interface for a component deciding
// the amount of instances needed of a job at a certain point in time.
type CapacityPlanner interface {
	// Plan plans how many instances are needed based on the given
	// scaleFactor
	Plan(scaleFactor float32, currentScale uint) uint

	// IsCoolingDown returns true if the CapacityPlanner thinks that
	// its currently not a good idea to apply the wanted scaling event.
	IsCoolingDown(timeOfLastScale time.Time, scaleFactor float32) bool
}
