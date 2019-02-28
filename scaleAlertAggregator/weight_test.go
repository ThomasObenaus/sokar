package scaleAlertAggregator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GetWeight(t *testing.T) {

	wm := map[string]float32{"AlertA": 2.0, "AlertB": -1}
	w := getWeight("AlertA", wm)
	assert.Equal(t, float32(2.0), w)

	w = getWeight("AlertB", wm)
	assert.Equal(t, float32(-1.0), w)

	w = getWeight("AlertC", wm)
	assert.Equal(t, float32(0), w)
}

func Test_weightPerSecondToWeight(t *testing.T) {

	assert.Equal(t, float32(1), weightPerSecondToWeight(1, time.Second))
	assert.Equal(t, float32(2), weightPerSecondToWeight(1, time.Second*2))
	assert.Equal(t, float32(0), weightPerSecondToWeight(1, time.Second*0))
	assert.Equal(t, float32(10), weightPerSecondToWeight(5, time.Second*2))
	assert.Equal(t, float32(-10), weightPerSecondToWeight(-5, time.Second*2))
	assert.Equal(t, float32(0), weightPerSecondToWeight(0, time.Second*2))
	assert.Equal(t, float32(0.5), weightPerSecondToWeight(1, time.Millisecond*500))
}
