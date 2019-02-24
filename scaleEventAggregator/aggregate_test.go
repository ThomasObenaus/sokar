package scaleEventAggregator

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

func Test_ComputeScaleValue(t *testing.T) {

	assert.Equal(t, float32(1), computeScaleValue(1, time.Second))
	assert.Equal(t, float32(2), computeScaleValue(1, time.Second*2))
	assert.Equal(t, float32(0), computeScaleValue(1, time.Second*0))
	assert.Equal(t, float32(10), computeScaleValue(5, time.Second*2))
	assert.Equal(t, float32(-10), computeScaleValue(-5, time.Second*2))
	assert.Equal(t, float32(0), computeScaleValue(0, time.Second*2))
	assert.Equal(t, float32(0.5), computeScaleValue(1, time.Millisecond*500))
}

func Test_ComputeScaleCounterDamping(t *testing.T) {
	assert.Equal(t, float32(0), computeScaleCounterDamping(0, 1))
	assert.Equal(t, float32(-1), computeScaleCounterDamping(1, 1))
	assert.Equal(t, float32(1), computeScaleCounterDamping(-1, 1))
	assert.Equal(t, float32(1), computeScaleCounterDamping(-10, 1))
	assert.Equal(t, float32(-1), computeScaleCounterDamping(10, 1))
}
