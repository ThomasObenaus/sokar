package sokar

// ScaleState represents the state of a scaling
type ScaleState string

const (
	// ScaleDone means the scale process was completed successfully
	ScaleDone ScaleState = "done"
	// ScaleRunning means the scale process is in progress
	ScaleRunning ScaleState = "running"
	// ScaleFailed means the scale process was completed but failed
	ScaleFailed ScaleState = "failed"
	// ScaleIgnored means the scale process was ignored (eventually not needed)
	ScaleIgnored ScaleState = "ignored"
)

type ScaleResult struct {
	State            ScaleState
	StateDescription string
	NewCount         uint
}

type Scaler interface {
	ScaleBy(amount int) ScaleResult
}
