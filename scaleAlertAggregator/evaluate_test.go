package scaleAlertAggregator

import (
	"testing"
	"time"

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

func Test_Evaluate(t *testing.T) {
	cfg := Config{}
	var emitters []ScaleAlertEmitter
	saa := cfg.New(emitters)

	saa.evaluationCycle = time.Second * 10
	saa.evaluationPeriodFactor = 10
	saa.downScalingThreshold = -10
	saa.upScalingThreshold = 10
	ago10Secs := time.Now().Add(time.Second * -10)

	// No Scale
	saa.evaluationCounter = 0
	saa.scaleCounter = 0
	saa.scaleCounterGradient.Update(saa.scaleCounter, ago10Secs)
	assert.Equal(t, float32(0), saa.evaluate())

	//  Scale UP
	saa.scaleCounterGradient.Update(0, ago10Secs)
	saa.evaluationCounter = 0
	saa.scaleCounter = 20
	assert.InDelta(t, float32(2), saa.evaluate(), 0.01)

	//  Scale DOWN
	saa.scaleCounterGradient.Update(0, ago10Secs)
	saa.evaluationCounter = 0
	saa.scaleCounter = -20
	assert.InDelta(t, float32(-2), saa.evaluate(), 0.01)
}
