package sokar

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"
	mock_sokar "github.com/thomasobenaus/sokar/test/sokar"
)

func Test_New(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)
	metrics, _ := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	assert.NotNil(t, sokar.scaleEventEmitter)
	assert.NotNil(t, sokar.capacityPlanner)
	assert.NotNil(t, sokar.scaler)
	assert.NotNil(t, sokar.schedule)
	assert.NotNil(t, sokar.stopChan)
	assert.NotNil(t, sokar.metrics)
}

func Test_HandleScaleEvent(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	scaleTo := uint(1)
	currentScale := uint(0)
	scaleFactor := float32(1)
	event := sokarIF.ScaleEvent{ScaleFactor: scaleFactor}
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		scalerIF.EXPECT().GetTimeOfLastScaleAction().Return(time.Now()),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), false).Return(false, time.Second*0),
		capaPlannerIF.EXPECT().Plan(scaleFactor, uint(0)).Return(scaleTo),
		scalerIF.EXPECT().ScaleTo(scaleTo, false),
	)
	metricMocks.scaleEventsTotal.EXPECT().Inc().Times(1)
	metricMocks.scaleFactor.EXPECT().Set(float64(scaleFactor))
	metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale))
	metricMocks.plannedJobCount.EXPECT().Set(float64(scaleTo))

	sokar.handleScaleEvent(event)
}

func Test_Run(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)
	metrics, _ := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	evEmitterIF.EXPECT().Subscribe(gomock.Any())
	sokar.Run()
}

func Test_ScaleValueToScaleDir(t *testing.T) {
	assert.True(t, scaleValueToScaleDir(-1))
	assert.False(t, scaleValueToScaleDir(1))
	assert.False(t, scaleValueToScaleDir(0))
}

func Test_TriggerScale_Scale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	currentScale := uint(0)
	scaleFactor := float32(1)
	scaleTo := uint(1)
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale)),
		scalerIF.EXPECT().GetTimeOfLastScaleAction().Return(time.Now()),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), false).Return(false, time.Second*0),
		metricMocks.plannedJobCount.EXPECT().Set(float64(scaleTo)),
		scalerIF.EXPECT().ScaleTo(scaleTo, false).Return(nil),
	)

	planFunc := func(scaleValue float32, currentScale uint) uint {
		return scaleTo
	}

	sokar.triggerScale(false, scaleFactor, planFunc)
}

func Test_TriggerScale_Cooldown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	currentScale := uint(0)
	scaleFactor := float32(1)
	scaleTo := uint(1)
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale)),
		scalerIF.EXPECT().GetTimeOfLastScaleAction().Return(time.Now()),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), false).Return(true, time.Second*0),
		metricMocks.skippedScalingDuringCooldownTotal.EXPECT().Inc(),
	)

	planFunc := func(scaleValue float32, currentScale uint) uint {
		return scaleTo
	}

	sokar.triggerScale(false, scaleFactor, planFunc)
}

func Test_TriggerScale_NoCooldown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	currentScale := uint(0)
	scaleFactor := float32(1)
	scaleTo := uint(1)
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale)),
		scalerIF.EXPECT().GetTimeOfLastScaleAction().Return(time.Now().Add(time.Hour*-1)),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), false).Return(false, time.Second*0),
		metricMocks.plannedJobCount.EXPECT().Set(float64(1)),
		scalerIF.EXPECT().ScaleTo(scaleTo, false),
	)

	planFunc := func(scaleValue float32, currentScale uint) uint {
		return scaleTo
	}

	sokar.triggerScale(false, scaleFactor, planFunc)
}

func Test_TriggerScale_ErrGettingJobCount(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	currentScale := uint(0)
	scaleFactor := float32(1)
	scaleTo := uint(1)
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, fmt.Errorf("Unable to obtain count")),
		metricMocks.failedScalingTotal.EXPECT().Inc(),
	)

	planFunc := func(scaleValue float32, currentScale uint) uint {
		return scaleTo
	}

	sokar.triggerScale(false, scaleFactor, planFunc)
}

func Test_TriggerScale_ErrScaleTo(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	scheduleIF := mock_sokar.NewMockScaleSchedule(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, scheduleIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	currentScale := uint(0)
	scaleFactor := float32(1)
	scaleTo := uint(1)
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale)),
		scalerIF.EXPECT().GetTimeOfLastScaleAction().Return(time.Now()),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), false).Return(false, time.Second*0),
		metricMocks.plannedJobCount.EXPECT().Set(float64(scaleTo)),
		scalerIF.EXPECT().ScaleTo(scaleTo, false).Return(fmt.Errorf("Unable to scale")),
		metricMocks.failedScalingTotal.EXPECT().Inc(),
	)

	planFunc := func(scaleValue float32, currentScale uint) uint {
		return scaleTo
	}

	sokar.triggerScale(false, scaleFactor, planFunc)
}
