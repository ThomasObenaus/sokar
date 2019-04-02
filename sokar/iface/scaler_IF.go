package sokar

// Scaler is a component that is able to scale a job/instance
type Scaler interface {
	ScaleTo(count uint, dryRun bool) error
	GetCount() (uint, error)
}
