package scaleEventAggregator

import (
	"testing"

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
