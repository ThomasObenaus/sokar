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
	metrics, _ := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	assert.NotNil(t, sokar.scaleEventEmitter)
	assert.NotNil(t, sokar.capacityPlanner)
	assert.NotNil(t, sokar.scaler)
	assert.NotNil(t, sokar.stopChan)
	assert.NotNil(t, sokar.metrics)

	oneDayAgo := time.Now().Add(time.Hour * -24)
	assert.WithinDuration(t, oneDayAgo, sokar.lastScaleAction, time.Second*1)
}

func Test_HandleScaleEvent(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	scaleTo := uint(1)
	currentScale := uint(0)
	scaleFactor := float32(1)
	event := sokarIF.ScaleEvent{ScaleFactor: scaleFactor}
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), scaleFactor).Return(false),
		capaPlannerIF.EXPECT().Plan(scaleFactor, uint(0)).Return(scaleTo),
		scalerIF.EXPECT().ScaleTo(scaleTo),
	)
	metricMocks.scaleEventsTotal.EXPECT().Inc().Times(1)
	metricMocks.scaleFactor.EXPECT().Set(float64(scaleFactor))
	metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale))
	metricMocks.plannedJobCount.EXPECT().Set(float64(scaleTo))

	sokar.handleScaleEvent(event)
}

func Test_HandleScaleEventOnCooldown(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	currentScale := uint(0)
	scaleFactor := float32(1)
	event := sokarIF.ScaleEvent{ScaleFactor: scaleFactor}
	gomock.InOrder(
		metricMocks.scaleEventsTotal.EXPECT().Inc().Times(1),
		metricMocks.scaleFactor.EXPECT().Set(float64(scaleFactor)),
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale)),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), scaleFactor).Return(true),
		metricMocks.skippedScalingDuringCooldownTotal.EXPECT().Inc(),
	)
	sokar.handleScaleEvent(event)
}

func Test_HandleScaleEvent_Fail(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	currentScale := uint(0)
	scaleFactor := float32(1)
	event := sokarIF.ScaleEvent{ScaleFactor: scaleFactor}
	scalerIF.EXPECT().GetCount().Return(currentScale, fmt.Errorf("Unable to obtain count"))
	metricMocks.scaleEventsTotal.EXPECT().Inc()
	metricMocks.scaleFactor.EXPECT().Set(float64(scaleFactor))
	metricMocks.failedScalingTotal.EXPECT().Inc()

	sokar.handleScaleEvent(event)

	scaleTo := uint(1)
	gomock.InOrder(
		metricMocks.scaleEventsTotal.EXPECT().Inc(),
		metricMocks.scaleFactor.EXPECT().Set(float64(scaleFactor)),
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale)),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), scaleFactor).Return(false),
		capaPlannerIF.EXPECT().Plan(scaleFactor, uint(0)).Return(scaleTo),
		metricMocks.plannedJobCount.EXPECT().Set(float64(scaleTo)),
		scalerIF.EXPECT().ScaleTo(scaleTo).Return(fmt.Errorf("ERROR")),
		metricMocks.failedScalingTotal.EXPECT().Inc(),
	)
	sokar.handleScaleEvent(event)
}

func Test_Run(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, _ := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	evEmitterIF.EXPECT().Subscribe(gomock.Any())
	sokar.Run()
}
