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

	lastScale = time.Now().Add(time.Second * -11)
	result = capa.IsCoolingDown(lastScale, false)
	assert.False(t, result)

	lastScale = time.Now().Add(time.Second * -21)
	result = capa.IsCoolingDown(lastScale, true)
	assert.False(t, result)
}
