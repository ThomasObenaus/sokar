package scaleAlertAggregator

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/thomasobenaus/sokar/test/metrics"
)

func Test_Metric(t *testing.T) {
	cfg := Config{}
	var emitters []ScaleAlertEmitter
	saa := cfg.New(emitters)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mScaleCounter := mock_metrics.NewMockGauge(mockCtrl)
	saa.metrics.scaleCounter = mScaleCounter
	mScaleCounter.EXPECT().Set(2)

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
