package helper

import (
	"sync"
	"time"
)

// LatestGradient is a structure that can be used to calculate the gradient of the two latest
// values of a field.
type LatestGradient struct {
	Value     float32
	Timestamp time.Time

	lock sync.Mutex
}

// UpdateAndGet is a thread safe method to calculate the gradient between
// the latest value and the current value.
// With each call to UpdateAndGet the given value and timestamp are stored internally.
// Furthermore the gradient is calculated by (currentValue-lastValue)/(currentTimeStamp-lastTimestamp).
func (lg *LatestGradient) UpdateAndGet(value float32, timestamp time.Time) float32 {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	increment := value - lg.Value
	timeSpan := timestamp.Sub(lg.Timestamp)

	lg.Value = value
	lg.Timestamp = timestamp

	if increment == 0 {
		return 0
	}

	if timeSpan.Seconds() <= 0 {
		return 0
	}

	return float32(float64(increment) / timeSpan.Seconds())
}

// Update sets new values for the gradient.
func (lg *LatestGradient) Update(value float32, timestamp time.Time) {
	lg.lock.Lock()
	defer lg.lock.Unlock()

	lg.Value = value
	lg.Timestamp = timestamp
}
