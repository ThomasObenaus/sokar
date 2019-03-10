package scaler

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/test/scaler"
)

func Test_New(t *testing.T) {

	cfg := Config{}
	scaler, err := cfg.New(nil)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg = Config{}
	scaler, err = cfg.New(scaTgt)
	assert.NoError(t, err)
	assert.NotNil(t, scaler)
	assert.NotNil(t, scaler.stopChan)
	assert.NotNil(t, scaler.scaleTicketChan)
	assert.NotNil(t, scaler.scalingTarget)
}

func Test_GetCount(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	scaTgt.EXPECT().GetJobCount("any").Return(uint(10), nil)

	cfg := Config{JobName: "any"}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	count, err := scaler.GetCount()
	assert.NoError(t, err)
	assert.Equal(t, uint(10), count)

	scaTgt.EXPECT().GetJobCount("any").Return(uint(0), fmt.Errorf("ERROR"))
	count, err = scaler.GetCount()
	assert.Error(t, err)
	assert.Equal(t, uint(0), count)
}

func Test_RunJoinStop(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg := Config{}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// Ensure that the run, join stop does not block forever
	scaler.Run()
	time.Sleep(time.Millisecond * 100)
	go func() {
		time.Sleep(time.Millisecond * 100)
		scaler.Stop()
	}()

	scaler.Join()
}
