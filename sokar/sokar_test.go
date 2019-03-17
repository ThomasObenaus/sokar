package sokar

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"
	"github.com/thomasobenaus/sokar/test/sokar"
)

func Test_HandleScaleEvent(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	scaleTo := uint(1)
	currentScale := uint(0)
	scaleFactor := float32(1)
	event := sokarIF.ScaleEvent{ScaleFactor: scaleFactor}
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		capaPlannerIF.EXPECT().Plan(scaleFactor, uint(0)).Return(scaleTo),
		scalerIF.EXPECT().ScaleTo(scaleTo),
	)
	sokar.handleScaleEvent(event)
}

func Test_Run(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	evEmitterIF.EXPECT().Subscribe(gomock.Any())
	sokar.Run()
}
