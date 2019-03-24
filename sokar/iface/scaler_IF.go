package sokar

// Scaler is a component that is able to scale a job/instance
type Scaler interface {
	ScaleTo(count uint) error
	GetCount() (uint, error)
}
