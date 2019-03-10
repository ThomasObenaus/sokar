package scaleAlertAggregator

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ComputeScaleCounterDamping(t *testing.T) {
	assert.Equal(t, float32(0), computeScaleCounterDamping(0, 1))
	assert.Equal(t, float32(-1), computeScaleCounterDamping(1, 1))
	assert.Equal(t, float32(1), computeScaleCounterDamping(-1, 1))
	assert.Equal(t, float32(1), computeScaleCounterDamping(-10, 1))
	assert.Equal(t, float32(-1), computeScaleCounterDamping(10, 1))
}

func Test_ApplyAlertsToScaleCounter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	var emitters []ScaleAlertEmitter
	saa := cfg.New(emitters, metrics)

	entries := make([]ScaleAlertPoolEntry, 0)
	entries = append(entries, ScaleAlertPoolEntry{weight: 2, scaleAlert: ScaleAlert{Name: "AlertA", Firing: true}})
	entries = append(entries, ScaleAlertPoolEntry{weight: -1, scaleAlert: ScaleAlert{Name: "AlertB", Firing: true}})

	saa.scaleCounter = 0
	counterHasChanged := saa.applyAlertsToScaleCounter(entries, time.Second*1)
	assert.True(t, counterHasChanged)
	assert.Equal(t, float32(1), saa.scaleCounter)

	saa.scaleCounter = 0
	entries = make([]ScaleAlertPoolEntry, 0)
	entries = append(entries, ScaleAlertPoolEntry{weight: 1, scaleAlert: ScaleAlert{Name: "AlertA", Firing: true}})
	entries = append(entries, ScaleAlertPoolEntry{weight: -1, scaleAlert: ScaleAlert{Name: "AlertB", Firing: true}})
	counterHasChanged = saa.applyAlertsToScaleCounter(entries, time.Second*1)
	assert.False(t, counterHasChanged)
	assert.Equal(t, float32(0), saa.scaleCounter)

	saa.scaleCounter = 0
	entries = make([]ScaleAlertPoolEntry, 0)
	entries = append(entries, ScaleAlertPoolEntry{weight: 1, scaleAlert: ScaleAlert{Name: "AlertA", Firing: true}})
	entries = append(entries, ScaleAlertPoolEntry{weight: -1, scaleAlert: ScaleAlert{Name: "AlertC", Firing: false}})
	counterHasChanged = saa.applyAlertsToScaleCounter(entries, time.Second*1)
	assert.True(t, counterHasChanged)
	assert.Equal(t, float32(1), saa.scaleCounter)
}

func Test_ApplyScaleCounterDamping(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	var emitters []ScaleAlertEmitter
	saa := cfg.New(emitters, metrics)

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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	gomock.InOrder(
		mocks.scaleCounter.EXPECT().Set(float64(0)),
		mocks.scaleCounter.EXPECT().Set(float64(1)),
		mocks.scaleCounter.EXPECT().Set(float64(-1)),
	)

	cfg := NewDefaultConfig()
	var emitters []ScaleAlertEmitter
	saa := cfg.New(emitters, metrics)
	saa.evaluationCycle = time.Second * 1

	// add some alerts to the pool
	alerts := make([]ScaleAlert, 0)
	alerts = append(alerts, ScaleAlert{Firing: true, Name: "AlertA"})
	alerts = append(alerts, ScaleAlert{Firing: true, Name: "AlertB"})

	weightMap := make(ScaleAlertWeightMap, 0)

	// No Scaling
	saa.scaleCounter = 0
	weightMap["AlertA"] = 1
	weightMap["AlertB"] = -1
	saa.scaleAlertPool.update("AM-Test", alerts, weightMap)
	require.Len(t, saa.scaleAlertPool.entries, 2)

	saa.aggregate()
	assert.Equal(t, float32(0), saa.scaleCounter)

	// Scaling Up
	saa.scaleCounter = 0
	weightMap["AlertA"] = 2
	weightMap["AlertB"] = -1
	saa.scaleAlertPool.update("AM-Test", alerts, weightMap)
	require.Len(t, saa.scaleAlertPool.entries, 2)

	saa.aggregate()
	assert.Equal(t, float32(1), saa.scaleCounter)

	// Scaling Down
	saa.scaleCounter = 0
	weightMap["AlertA"] = 1
	weightMap["AlertB"] = -2
	saa.scaleAlertPool.update("AM-Test", alerts, weightMap)
	require.Len(t, saa.scaleAlertPool.entries, 2)
	saa.aggregate()
	assert.Equal(t, float32(-1), saa.scaleCounter)
}
