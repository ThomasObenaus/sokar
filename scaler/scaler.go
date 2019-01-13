package scaler

// Scaler specifies the interface for a component that can scale a certain job
type Scaler interface {
	// ScaleBy
	ScaleBy(jobName string, amount int) error
}
