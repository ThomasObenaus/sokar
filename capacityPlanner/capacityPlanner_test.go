package capacityPlanner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	cfg := NewDefaultConfig()
	capa := cfg.New()
	assert.NotNil(t, capa)
}

func Test_Plan_ModeLinear(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.Mode = CapaPlanningModeLinear
	capa := cfg.New()
	require.NotNil(t, capa)

	assert.Equal(t, uint(10), capa.Plan(0, 10))

	assert.Equal(t, uint(1), capa.Plan(1, 0))
	assert.Equal(t, uint(20), capa.Plan(1, 10))
	assert.Equal(t, uint(2), capa.Plan(0.5, 1))
	assert.Equal(t, uint(15), capa.Plan(0.5, 10))

	assert.Equal(t, uint(0), capa.Plan(-1, 0))
	assert.Equal(t, uint(0), capa.Plan(-1, 10))
	assert.Equal(t, uint(0), capa.Plan(-0.5, 1))
	assert.Equal(t, uint(5), capa.Plan(-0.5, 10))
}

func Test_Plan_ModeConstant(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.Mode = CapaPlanningModeConstant

	cfg.OffsetConstantMode = 1
	capa := cfg.New()
	require.NotNil(t, capa)
	assert.Equal(t, uint(0), capa.Plan(-1, 0))
	assert.Equal(t, uint(0), capa.Plan(-1, 1))
	assert.Equal(t, uint(0), capa.Plan(0, 0))
	assert.Equal(t, uint(2), capa.Plan(0, 2))
	assert.Equal(t, uint(1), capa.Plan(1, 0))
	assert.Equal(t, uint(2), capa.Plan(1, 1))

	cfg.OffsetConstantMode = 2
	capa = cfg.New()
	require.NotNil(t, capa)
	assert.Equal(t, uint(0), capa.Plan(-1, 1))
	assert.Equal(t, uint(3), capa.Plan(-1, 5))
	assert.Equal(t, uint(3), capa.Plan(1, 1))
	assert.Equal(t, uint(7), capa.Plan(1, 5))

	cfg.OffsetConstantMode = 0
	capa = cfg.New()
	require.NotNil(t, capa)
	assert.Equal(t, uint(0), capa.Plan(1, 0))

}
func Test_IsCoolingDown(t *testing.T) {
	downScalePeriod := time.Second * 20
	upScalePeriod := time.Second * 10
	cfg := Config{
		DownScaleCooldownPeriod: downScalePeriod,
		UpScaleCooldownPeriod:   upScalePeriod,
	}
	capa := cfg.New()
	require.NotNil(t, capa)

	lastScale := time.Now()
	result := capa.IsCoolingDown(lastScale, false)
	assert.True(t, result)

	result = capa.IsCoolingDown(lastScale, true)
	assert.True(t, result)

	// Upscaling
	lastScale = time.Now().Add(time.Second * -11)
	result = capa.IsCoolingDown(lastScale, false)
	assert.False(t, result)

	lastScale = time.Now().Add(time.Second * -9)
	result = capa.IsCoolingDown(lastScale, false)
	assert.True(t, result)

	// Downscaling
	lastScale = time.Now().Add(time.Second * -21)
	result = capa.IsCoolingDown(lastScale, true)
	assert.False(t, result)

	lastScale = time.Now().Add(time.Second * -19)
	result = capa.IsCoolingDown(lastScale, true)
	assert.True(t, result)
}
