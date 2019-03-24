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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	scaler, err := cfg.New(nil, metrics)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg = Config{}
	scaler, err = cfg.New(scaTgt, metrics)
	assert.NoError(t, err)
	assert.NotNil(t, scaler)
	assert.NotNil(t, scaler.stopChan)
	assert.NotNil(t, scaler.scaleTicketChan)
	assert.NotNil(t, scaler.scalingTarget)
}

func Test_GetCount(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	scaTgt.EXPECT().GetJobCount("any").Return(uint(10), nil)

	cfg := Config{JobName: "any"}
	scaler, err := cfg.New(scaTgt, metrics)
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
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg := Config{}
	scaler, err := cfg.New(scaTgt, metrics)
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

func Test_OpenScalingTicket(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg := Config{}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	err = scaler.openScalingTicket(0)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), scaler.numOpenScalingTickets)
	assert.Len(t, scaler.scaleTicketChan, 1)

	err = scaler.openScalingTicket(0)
	assert.Error(t, err)

	ticket := <-scaler.scaleTicketChan
	assert.Equal(t, uint(0), ticket.desiredCount)
}

func Test_ApplyScalingTicket(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	scaTgt.EXPECT().GetJobCount("any").Return(uint(0), nil)
	scaTgt.EXPECT().IsJobDead("any").Return(true, nil)

	cfg := Config{JobName: "any"}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	ticket := NewScalingTicket(0)
	scaler.applyScaleTicket(ticket)
}

func Test_OpenAndApplyScalingTicket(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	scaTgt.EXPECT().GetJobCount("any").Return(uint(0), nil).MaxTimes(11)
	scaTgt.EXPECT().IsJobDead("any").Return(true, nil).MaxTimes(11)

	cfg := Config{JobName: "any", MaxOpenScalingTickets: 10}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// open as many tickets as allowed
	for i := uint(0); i <= scaler.maxOpenScalingTickets; i++ {
		err = scaler.openScalingTicket(0)
		assert.NoError(t, err)
	}

	// open new ticket --> should fail
	err = scaler.openScalingTicket(0)
	assert.Error(t, err)

	// apply/ close as many tickets as are open
	for i := uint(0); i <= scaler.maxOpenScalingTickets; i++ {
		ticket := <-scaler.scaleTicketChan
		scaler.applyScaleTicket(ticket)
	}

	// open new ticket --> should NOT fail
	err = scaler.openScalingTicket(0)
	assert.NoError(t, err)
}
