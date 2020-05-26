package main

import (
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	mock_main "github.com/thomasobenaus/sokar/test/mocks"
)

func Test_Run(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zerolog.Logger{}

	var runnables []Runnable
	runnable1 := mock_main.NewMockRunnable(mockCtrl)
	runnable2 := mock_main.NewMockRunnable(mockCtrl)

	gomock.InOrder(
		runnable1.EXPECT().String().Times(1),
		runnable1.EXPECT().Start().Times(1),
		runnable2.EXPECT().String().Times(1),
		runnable2.EXPECT().Start().Times(1),
	)

	runnables = append(runnables, runnable1)
	runnables = append(runnables, runnable2)
	Run(runnables, logger)
}

func Test_Join(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zerolog.Logger{}

	var runnables []Runnable
	runnable1 := mock_main.NewMockRunnable(mockCtrl)
	runnable2 := mock_main.NewMockRunnable(mockCtrl)

	gomock.InOrder(
		runnable1.EXPECT().String().Times(1),
		runnable1.EXPECT().Join().Times(1),
		runnable2.EXPECT().String().Times(1),
		runnable2.EXPECT().Join().Times(1),
	)

	runnables = append(runnables, runnable1)
	runnables = append(runnables, runnable2)
	Join(runnables, logger)
}

func Test_Stop(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zerolog.Logger{}

	var runnables []Runnable
	runnable1 := mock_main.NewMockRunnable(mockCtrl)
	runnable2 := mock_main.NewMockRunnable(mockCtrl)

	gomock.InOrder(
		runnable2.EXPECT().String().Times(1),
		runnable2.EXPECT().Stop().Times(1),
		runnable1.EXPECT().String().Times(1),
		runnable1.EXPECT().Stop().Times(1),
	)

	runnables = append(runnables, runnable1)
	runnables = append(runnables, runnable2)
	Stop(runnables, logger)
}

type testSignal struct {
}

func (s testSignal) String() string {
	return "testSignal"
}

func (s testSignal) Signal() {
}

func Test_Shutdown(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zerolog.Logger{}

	var runnables []Runnable
	runnable1 := mock_main.NewMockRunnable(mockCtrl)
	runnable2 := mock_main.NewMockRunnable(mockCtrl)

	gomock.InOrder(
		runnable2.EXPECT().String().Times(1),
		runnable2.EXPECT().Stop().Times(1),
		runnable1.EXPECT().String().Times(1),
		runnable1.EXPECT().Stop().Times(1),
	)

	runnables = append(runnables, runnable1)
	runnables = append(runnables, runnable2)

	shutDownChan := make(chan os.Signal, 1)
	go shutdownHandler(shutDownChan, runnables, logger)
	time.Sleep(time.Millisecond * 100)

	shutDownChan <- testSignal{}
	time.Sleep(time.Millisecond * 100)
}
