package scaler

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/test/metrics"
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
	scaTgt.EXPECT().GetScalingObjectCount("any").Return(uint(10), nil)

	cfg := Config{Name: "any"}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	count, err := scaler.GetCount()
	assert.NoError(t, err)
	assert.Equal(t, uint(10), count)

	scaTgt.EXPECT().GetScalingObjectCount("any").Return(uint(0), fmt.Errorf("ERROR"))
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

func Test_UpdateScaleResultMetric(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	scaleResultCounter := mock_metrics.NewMockCounterVec(mockCtrl)

	otherCounter := mock_metrics.NewMockCounter(mockCtrl)
	otherCounter.EXPECT().Inc()
	scaleResultCounter.EXPECT().WithLabelValues("other").Return(otherCounter)
	updateScaleResultMetric(scaleResult{}, scaleResultCounter)

	failedCounter := mock_metrics.NewMockCounter(mockCtrl)
	failedCounter.EXPECT().Inc()
	scaleResultCounter.EXPECT().WithLabelValues("failed").Return(failedCounter)
	updateScaleResultMetric(scaleResult{state: scaleFailed}, scaleResultCounter)

	doneCounter := mock_metrics.NewMockCounter(mockCtrl)
	doneCounter.EXPECT().Inc()
	scaleResultCounter.EXPECT().WithLabelValues("done").Return(doneCounter)
	updateScaleResultMetric(scaleResult{state: scaleDone}, scaleResultCounter)

	ignoredCounter := mock_metrics.NewMockCounter(mockCtrl)
	ignoredCounter.EXPECT().Inc()
	scaleResultCounter.EXPECT().WithLabelValues("ignored").Return(ignoredCounter)
	updateScaleResultMetric(scaleResult{state: scaleIgnored}, scaleResultCounter)
}

func Test_OpenScalingTicket(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	cfg := Config{}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc().Times(2)
	mocks.scalingTicketCount.EXPECT().WithLabelValues("added").Return(scalingTicketCounter)
	err = scaler.openScalingTicket(0, false)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), scaler.numOpenScalingTickets)
	assert.Len(t, scaler.scaleTicketChan, 1)

	mocks.scalingTicketCount.EXPECT().WithLabelValues("rejected").Return(scalingTicketCounter)
	err = scaler.openScalingTicket(0, false)
	assert.Error(t, err)

	ticket := <-scaler.scaleTicketChan
	assert.Equal(t, uint(0), ticket.desiredCount)
}

func Test_ApplyScalingTicket(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	scaTgt.EXPECT().GetScalingObjectCount("any").Return(uint(0), nil)
	scaTgt.EXPECT().IsScalingObjectDead("any").Return(true, nil)

	cfg := Config{Name: "any"}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc()
	mocks.scalingTicketCount.EXPECT().WithLabelValues("applied").Return(scalingTicketCounter)
	ignoredCounter := mock_metrics.NewMockCounter(mockCtrl)
	ignoredCounter.EXPECT().Inc()
	mocks.scaleResultCounter.EXPECT().WithLabelValues("ignored").Return(ignoredCounter)
	mocks.scalingDurationSeconds.EXPECT().Observe(gomock.Any())

	ticket := NewScalingTicket(0, false)
	scaler.applyScaleTicket(ticket)
}

func Test_OpenAndApplyScalingTicket(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	scaTgt.EXPECT().GetScalingObjectCount("any").Return(uint(0), nil).MaxTimes(11)
	scaTgt.EXPECT().IsScalingObjectDead("any").Return(true, nil).MaxTimes(11)

	cfg := Config{Name: "any", MaxOpenScalingTickets: 10}
	scaler, err := cfg.New(scaTgt, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// open as many tickets as allowed
	ticketCounter := int(scaler.maxOpenScalingTickets) + 1
	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc().Times(ticketCounter)
	mocks.scalingTicketCount.EXPECT().WithLabelValues("added").Return(scalingTicketCounter).Times(ticketCounter)

	for i := uint(0); i <= scaler.maxOpenScalingTickets; i++ {
		err = scaler.openScalingTicket(0, false)
		assert.NoError(t, err)
	}

	// open new ticket --> should fail
	scalingTicketCounter.EXPECT().Inc().Times(1)
	mocks.scalingTicketCount.EXPECT().WithLabelValues("rejected").Return(scalingTicketCounter)
	err = scaler.openScalingTicket(0, false)
	assert.Error(t, err)

	// apply/ close as many tickets as are open
	scalingTicketCounter.EXPECT().Inc().Times(ticketCounter)
	mocks.scalingTicketCount.EXPECT().WithLabelValues("applied").Return(scalingTicketCounter).Times(ticketCounter)
	ignoredCounter := mock_metrics.NewMockCounter(mockCtrl)
	ignoredCounter.EXPECT().Inc().Times(ticketCounter)
	mocks.scaleResultCounter.EXPECT().WithLabelValues("ignored").Return(ignoredCounter).Times(ticketCounter)
	mocks.scalingDurationSeconds.EXPECT().Observe(gomock.Any()).Times(ticketCounter)
	for i := uint(0); i <= scaler.maxOpenScalingTickets; i++ {
		ticket := <-scaler.scaleTicketChan
		scaler.applyScaleTicket(ticket)
	}

	// open new ticket --> should NOT fail
	scalingTicketCounter.EXPECT().Inc().Times(1)
	mocks.scalingTicketCount.EXPECT().WithLabelValues("added").Return(scalingTicketCounter).Times(1)
	err = scaler.openScalingTicket(0, false)
	assert.NoError(t, err)
}
