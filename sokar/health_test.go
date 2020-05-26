package sokar

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_sokar "github.com/thomasobenaus/sokar/test/mocks/sokar"
)

func Test_IsHealthy(t *testing.T) {

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

	err = sokar.IsHealthy()
	assert.NoError(t, err)
}
