package scaler

// ScalingObject config of the ScalingObject to be scaled
type ScalingObject struct {
	Name     string
	MinCount uint
	MaxCount uint
}
