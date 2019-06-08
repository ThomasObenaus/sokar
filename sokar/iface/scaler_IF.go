package sokar

// Scaler is a component that is able to scale a scaling-object
type Scaler interface {
	ScaleTo(count uint, dryRun bool) error
	GetCount() (uint, error)
}
