package sokar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PlanScaleByPercentage(t *testing.T) {

	assert.Equal(t, uint(2), planScaleByPercentage(0.1, 1))
	assert.Equal(t, uint(0), planScaleByPercentage(-0.1, 1))
	assert.Equal(t, uint(110), planScaleByPercentage(0.1, 100))
	assert.Equal(t, uint(90), planScaleByPercentage(-0.1, 100))
	assert.Equal(t, uint(0), planScaleByPercentage(-1, 100))
	assert.Equal(t, uint(200), planScaleByPercentage(1, 100))
	assert.Equal(t, uint(300), planScaleByPercentage(2, 100))
	assert.Equal(t, uint(0), planScaleByPercentage(-1.1, 100))

	assert.Equal(t, uint(33), planScaleByPercentage(10, 3))
}
func Test_PlanScaleByValue(t *testing.T) {

	assert.Equal(t, uint(100), planScaleByValue(0.1, 100))
	assert.Equal(t, uint(101), planScaleByValue(1, 100))
	assert.Equal(t, uint(101), planScaleByValue(0.9, 100))
	assert.Equal(t, uint(200), planScaleByValue(100, 100))
	assert.Equal(t, uint(100), planScaleByValue(-0.1, 100))
	assert.Equal(t, uint(99), planScaleByValue(-1, 100))
	assert.Equal(t, uint(99), planScaleByValue(-0.9, 100))
	assert.Equal(t, uint(0), planScaleByValue(-100, 100))
	assert.Equal(t, uint(0), planScaleByValue(-101, 100))
}
