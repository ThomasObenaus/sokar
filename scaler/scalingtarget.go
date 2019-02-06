package scaler

// ScalingTarget represents the interface to be implemented
// in order to be used by the Scaler as scaling target.
type ScalingTarget interface {
	SetJobCount(jobname string, count uint) error
	GetJobCount(jobname string) (uint, error)
	GetJobState(jobname string) (JobState, error)
}

// JobState represents border types
type JobState uint

const (
	// JobStateUnknown if state couldn't be determined
	JobStateUnknown JobState = iota
	// JobStateRunning Job is running
	JobStateRunning
	// JobStatePending Job is pending
	JobStatePending
	// JobStateDead Job is dead
	JobStateDead
)
