package scaler

// ScalingTarget represents the interface to be implemented
// in order to be used by the Scaler as scaling target.
type ScalingTarget interface {
	SetJobCount(jobname string, count uint) error
	GetJobCount(jobname string) (uint, error)
}
