package scaleAlertAggregator

import (
	"fmt"
	"time"
)

type scaleCounter struct {
	firstTimeChanged     time.Time
	val                  float32
	wasChangedAfterReset bool
}

func newScaleCounter() scaleCounter {
	result := scaleCounter{}
	result.reset()
	return result
}

func (sc *scaleCounter) reset() {
	sc.val = 0
	sc.wasChangedAfterReset = false
	sc.firstTimeChanged = time.Unix(0, 0)
}

func (sc *scaleCounter) incBy(val float32) {
	if !sc.wasChangedAfterReset && val != 0 {
		sc.wasChangedAfterReset = true
		sc.firstTimeChanged = time.Now()
	}
	sc.val += val
}

func (sc *scaleCounter) durationSinceFirstChange() (time.Duration, error) {
	if sc.wasChangedAfterReset {
		return time.Now().Sub(sc.firstTimeChanged), nil
	}
	return 0, fmt.Errorf("Can calculate duration since the counter was never changed")
}
