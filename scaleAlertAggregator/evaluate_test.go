package scaleAlertAggregator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsScalingNeeded(t *testing.T) {
	cfg := Config{}
	var emitters []ScaleAlertEmitter
	saa := cfg.New(emitters)
	saa.downScalingThreshold = -5
	saa.upScalingThreshold = 5

	saa.scaleCounter = 0
	assert.False(t, saa.isScalingNeeded())

	saa.scaleCounter = 5.1
	assert.True(t, saa.isScalingNeeded())

	saa.scaleCounter = -5.1
	assert.True(t, saa.isScalingNeeded())
}

func Test_GradientToScaleDir(t *testing.T) {
	assert.Equal(t, "UP", gradientToScaleDir(1))
	assert.Equal(t, "DOWN", gradientToScaleDir(-1))
	assert.Equal(t, "NO", gradientToScaleDir(0))
}
