package scaler

// Scaler specifies the interface for a component that can scale a certain job
type Scaler interface {
	// ScaleBy Scales the target component by the given amount of instances
	ScaleBy(amount int) error
}
