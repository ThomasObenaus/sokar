package scaler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetOptionalValue(t *testing.T) {
	optVal := optionalValue{}
	assert.False(t, optVal.isKnown)
	optVal.setValue(1)
	assert.True(t, optVal.isKnown)
	assert.Equal(t, uint(1), optVal.value)
}
