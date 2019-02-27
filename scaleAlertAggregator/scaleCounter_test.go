package scaleAlertAggregator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	sc := newScaleCounter()

	assert.Equal(t, float32(0), sc.val)
	assert.Equal(t, time.Unix(0, 0), sc.firstTimeChanged)
	assert.False(t, sc.wasChangedAfterReset)
}

func Test_Reset(t *testing.T) {
	sc := newScaleCounter()
	sc.firstTimeChanged = time.Now()
	sc.val = 10
	sc.wasChangedAfterReset = true

	sc.reset()
	assert.Equal(t, float32(0), sc.val)
	assert.Equal(t, time.Unix(0, 0), sc.firstTimeChanged)
	assert.False(t, sc.wasChangedAfterReset)
}

func Test_IncBy(t *testing.T) {
	sc := newScaleCounter()
	sc.reset()

	ago10Sec := time.Now().Add(time.Second * -10)
	sc.incBy(10.3)

	assert.Equal(t, float32(10.3), sc.val)
	assert.True(t, sc.firstTimeChanged.After(ago10Sec))
	assert.True(t, sc.wasChangedAfterReset)

	changedAt := sc.firstTimeChanged
	sc.incBy(10)

	assert.Equal(t, float32(20.3), sc.val)
	assert.WithinDuration(t, changedAt, sc.firstTimeChanged, time.Millisecond*0)
}

func Test_DurationSinceFirstChange(t *testing.T) {
	sc := newScaleCounter()
	sc.reset()
	sc.incBy(1)
	time.Sleep(1 * time.Second)

	dur, err := sc.durationSinceFirstChange()
	assert.NoError(t, err)
	assert.InDelta(t, time.Second*1, dur, float64(time.Millisecond*1))

	sc.reset()
	dur, err = sc.durationSinceFirstChange()
	assert.Error(t, err)
}
