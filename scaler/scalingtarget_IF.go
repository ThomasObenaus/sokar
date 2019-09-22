package scaler

// ScalingTarget represents the interface to be implemented
// in order to be used by the Scaler as scaling target.
type ScalingTarget interface {
	SetScalingObjectCount(scalingObject string, count uint) error
	GetScalingObjectCount(scalingObject string) (uint, error)
	IsScalingObjectDead(scalingObject string) (bool, error)
	String() string
}
