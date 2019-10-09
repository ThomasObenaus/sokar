package sokar

import "time"

// Scaler is a component that is able to scale a scaling-object
type Scaler interface {
	ScaleTo(count uint, force bool) error
	GetCount() (uint, error)
	GetTimeOfLastScaleAction() time.Time
}
