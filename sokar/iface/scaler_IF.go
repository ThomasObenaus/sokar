package sokar

// ScaleState represents the state of a scaling
type ScaleState string

const (
	// ScaleUnknown means the scale process was completed successfully
	ScaleUnknown ScaleState = "unknown"
	// ScaleDone means the scale process was completed successfully
	ScaleDone ScaleState = "done"
	// ScaleRunning means the scale process is in progress
	ScaleRunning ScaleState = "running"
	// ScaleFailed means the scale process was completed but failed
	ScaleFailed ScaleState = "failed"
	// ScaleIgnored means the scale process was ignored (eventually not needed)
	ScaleIgnored ScaleState = "ignored"
	// ScaleNotStarted means the scale process was not started yet
	ScaleNotStarted ScaleState = "not started"
)

// ScaleResult is created after scaling was done and contains the result
type ScaleResult struct {
	State            ScaleState
	StateDescription string
	NewCount         uint
}

// Scaler is a component that is able to scale a job/instance
type Scaler interface {
	ScaleBy_Old(amount int) ScaleResult
	ScaleTo(count uint) ScaleResult
}
