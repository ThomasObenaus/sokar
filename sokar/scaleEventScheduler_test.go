package sokar

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_sokar "github.com/thomasobenaus/sokar/test/sokar"
)

func Test_ShouldFireScaleEvent(t *testing.T) {

	// GIVEN
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

	// WHEN
	scheduleIF.EXPECT().IsActiveAt(gomock.Any(), gomock.Any()).Return(true)
	result := sokar.shouldFireScaleEvent(time.Now())

	// THEN
	assert.True(t, result)
}

func Test_ShouldNotFireScaleEvent(t *testing.T) {

	// GIVEN
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

	// WHEN
	scheduleIF.EXPECT().IsActiveAt(gomock.Any(), gomock.Any()).Return(false)
	result := sokar.shouldFireScaleEvent(time.Now())

	// THEN
	assert.False(t, result)
}
