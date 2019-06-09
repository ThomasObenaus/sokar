package capacityPlanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PlanConstant(t *testing.T) {

	cfg := NewDefaultConfig()
	cap, err := cfg.New()
	require.NotNil(t, cap)
	require.NoError(t, err)

	assert.Equal(t, uint(0), cap.planConstant(-1, 0, 1))
	assert.Equal(t, uint(0), cap.planConstant(-1, 1, 1))
	assert.Equal(t, uint(0), cap.planConstant(-1, 1, 2))
	assert.Equal(t, uint(3), cap.planConstant(-1, 5, 2))

	assert.Equal(t, uint(0), cap.planConstant(1, 0, 0))
	assert.Equal(t, uint(1), cap.planConstant(1, 0, 1))
	assert.Equal(t, uint(2), cap.planConstant(1, 1, 1))
	assert.Equal(t, uint(3), cap.planConstant(1, 1, 2))
	assert.Equal(t, uint(7), cap.planConstant(1, 5, 2))

	assert.Equal(t, uint(0), cap.planConstant(0, 0, 1))
	assert.Equal(t, uint(2), cap.planConstant(0, 2, 1))
}
