package scaleAlertAggregator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_ComputeScaleCounterDamping(t *testing.T) {
	assert.Equal(t, float32(0), computeScaleCounterDamping(0, 1))
	assert.Equal(t, float32(-1), computeScaleCounterDamping(1, 1))
	assert.Equal(t, float32(1), computeScaleCounterDamping(-1, 1))
	assert.Equal(t, float32(1), computeScaleCounterDamping(-10, 1))
	assert.Equal(t, float32(-1), computeScaleCounterDamping(10, 1))
}

func Test_ComputeScaleCounterIncrement(t *testing.T) {

	wm := map[string]float32{"AlertA": 2.0, "AlertB": -1}

	inc, wps := computeScaleCounterIncrement("AlertA", wm, time.Second*1)
	assert.Equal(t, float32(2), inc)
	assert.Equal(t, float32(2), wps)

	inc, _ = computeScaleCounterIncrement("AlertB", wm, time.Second*1)
	assert.Equal(t, float32(-1), inc)

	inc, _ = computeScaleCounterIncrement("AlertB", wm, time.Second*2)
	assert.Equal(t, float32(-2), inc)

	inc, _ = computeScaleCounterIncrement("AlertA", wm, time.Second*2)
	assert.Equal(t, float32(4), inc)

	inc, _ = computeScaleCounterIncrement("AlertA", wm, time.Millisecond*500)
	assert.Equal(t, float32(1), inc)

	inc, _ = computeScaleCounterIncrement("AlertB", wm, time.Millisecond*500)
	assert.Equal(t, float32(-0.5), inc)

	inc, wps = computeScaleCounterIncrement("NO ALERT", wm, time.Second*2)
	assert.Equal(t, float32(0), inc)
	assert.Equal(t, float32(0), wps)
}

func Test_ApplyAlertsToScaleCounter(t *testing.T) {
	cfg := Config{}
	var emitters []ScaleAlertEmitter
	saa := cfg.New(emitters)

	var entries []ScaleAlertPoolEntry
	entries = append(entries, ScaleAlertPoolEntry{scaleAlert: ScaleAlert{Name: "AlertA", Firing: true}})
	entries = append(entries, ScaleAlertPoolEntry{scaleAlert: ScaleAlert{Name: "AlertB", Firing: true}})
	entries = append(entries, ScaleAlertPoolEntry{scaleAlert: ScaleAlert{Name: "AlertC", Firing: false}})

	saa.scaleCounter = 0
	wm := map[string]float32{"AlertA": 2.0, "AlertB": -1}
	counterHasChanged := saa.applyAlertsToScaleCounter(entries, wm, time.Second*1)
	assert.True(t, counterHasChanged)
	assert.Equal(t, float32(1), saa.scaleCounter)

	saa.scaleCounter = 0
	wm = map[string]float32{"AlertA": 1.0, "AlertB": -1}
	counterHasChanged = saa.applyAlertsToScaleCounter(entries, wm, time.Second*1)
	assert.False(t, counterHasChanged)
	assert.Equal(t, float32(0), saa.scaleCounter)

	saa.scaleCounter = 0
	wm = map[string]float32{"AlertA": 1.0, "AlertC": -1}
	counterHasChanged = saa.applyAlertsToScaleCounter(entries, wm, time.Second*1)
	assert.True(t, counterHasChanged)
	assert.Equal(t, float32(1), saa.scaleCounter)
}

func Test_ApplyScaleCounterDamping(t *testing.T) {
	cfg := Config{}
	var emitters []ScaleAlertEmitter
	saa := cfg.New(emitters)

	saa.scaleCounter = 0
	saa.applyScaleCounterDamping(1, time.Second*1)
	assert.Equal(t, float32(0), saa.scaleCounter)

	saa.scaleCounter = 2
	saa.applyScaleCounterDamping(1, time.Second*1)
	assert.Equal(t, float32(1), saa.scaleCounter)

	saa.scaleCounter = -2
	saa.applyScaleCounterDamping(1, time.Second*1)
	assert.Equal(t, float32(-1), saa.scaleCounter)
}

func Test_Aggregate(t *testing.T) {
	cfg := Config{}
	var emitters []ScaleAlertEmitter
	saa := cfg.New(emitters)

	// add some alerts to the pool
	var alerts []ScaleAlert
	alerts = append(alerts, ScaleAlert{Firing: true, Name: "AlertA"})
	alerts = append(alerts, ScaleAlert{Firing: true, Name: "AlertB"})
	alerts = append(alerts, ScaleAlert{Firing: false, Name: "AlertC"})
	saa.scaleAlertPool.update("AM-Test", alerts)
	require.Len(t, saa.scaleAlertPool.entries, 2)
	saa.upScalingThreshold = 1
	saa.downScalingThreshold = -1

	// No Scaling
	saa.scaleCounter = 0
	saa.weightMap["AlertA"] = 1
	saa.weightMap["AlertB"] = -1
	assert.Equal(t, float32(0), saa.scaleCounter)
	assert.Equal(t, float32(0), saa.aggregate())

	// Scaling Up
	saa.scaleCounter = 0
	saa.weightMap["AlertA"] = 2
	saa.weightMap["AlertB"] = -1
	assert.Equal(t, float32(0), saa.scaleCounter)
	assert.Equal(t, float32(1), saa.aggregate())

	// Scaling Down
	saa.scaleCounter = 0
	saa.weightMap["AlertA"] = 1
	saa.weightMap["AlertB"] = -2
	assert.Equal(t, float32(0), saa.scaleCounter)
	assert.Equal(t, float32(-1), saa.aggregate())
}
