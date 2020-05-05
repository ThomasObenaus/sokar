package scaler

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_metrics "github.com/thomasobenaus/sokar/test/metrics"
	mock_scaler "github.com/thomasobenaus/sokar/test/scaler"
)

func Test_New(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	sObj := ScalingObject{}
	scaler, err := New(nil, sObj, metrics)
	assert.Error(t, err)
	assert.Nil(t, scaler)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	scaler, err = New(scaTgt, sObj, metrics)
	assert.NoError(t, err)
	require.NotNil(t, scaler)
	assert.NotNil(t, scaler.stopChan)
	assert.NotNil(t, scaler.scaleTicketChan)
	assert.NotNil(t, scaler.scalingTarget)

	oneDayAgo := time.Now().Add(time.Hour * -24)
	assert.WithinDuration(t, oneDayAgo, scaler.lastScaleAction, time.Second*1)
}

func Test_GetCount(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	scaTgt.EXPECT().GetScalingObjectCount("any").Return(uint(10), nil)

	sObjName := "any"
	sObj := ScalingObject{Name: sObjName}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	count, err := scaler.GetCount()
	assert.NoError(t, err)
	assert.Equal(t, uint(10), count)

	scaTgt.EXPECT().GetScalingObjectCount(sObjName).Return(uint(0), fmt.Errorf("ERROR"))
	count, err = scaler.GetCount()
	assert.Error(t, err)
	assert.Equal(t, uint(0), count)
}

func Test_RunJoinStop(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	sObj := ScalingObject{}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	// Ensure that the run, join stop does not block forever
	scaler.Run()
	time.Sleep(time.Millisecond * 100)
	go func() {
		time.Sleep(time.Millisecond * 100)
		err := scaler.Stop()
		assert.NoError(t, err)
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

	sObj := ScalingObject{}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc().Times(2)
	mocks.scalingTicketCount.EXPECT().WithLabelValues("added").Return(scalingTicketCounter)
	err = scaler.openScalingTicket(1, false)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), scaler.numOpenScalingTickets)
	assert.Len(t, scaler.scaleTicketChan, 1)

	mocks.scalingTicketCount.EXPECT().WithLabelValues("rejected").Return(scalingTicketCounter)
	err = scaler.openScalingTicket(0, false)
	assert.Error(t, err)

	ticket := <-scaler.scaleTicketChan
	assert.Equal(t, uint(1), ticket.desiredCount)
	assert.False(t, ticket.force)
}

func Test_ApplyScalingTicket_NoScale_DeadJob(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	sObjName := "any"

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	scaTgt.EXPECT().GetScalingObjectCount(sObjName).Return(uint(0), nil)
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(true, nil)

	sObj := ScalingObject{Name: sObjName}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc()
	mocks.scalingTicketCount.EXPECT().WithLabelValues("applied").Return(scalingTicketCounter)
	ignoredCounter := mock_metrics.NewMockCounter(mockCtrl)
	ignoredCounter.EXPECT().Inc()
	mocks.scaleResultCounter.EXPECT().WithLabelValues("ignored").Return(ignoredCounter)
	mocks.scalingDurationSeconds.EXPECT().Observe(gomock.Any())

	ticket := NewScalingTicket(10, false)
	scaler.applyScaleTicket(ticket)
	assert.False(t, scaler.desiredScale.isKnown)
	assert.Equal(t, uint(0), scaler.desiredScale.value)
}

func Test_ApplyScalingTicket_NoScale_DryRun(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(1))
	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc()
	doneCounter := mock_metrics.NewMockCounter(mockCtrl)
	doneCounter.EXPECT().Inc()
	sObjName := "any"

	gomock.InOrder(
		scaTgt.EXPECT().GetScalingObjectCount(sObjName).Return(uint(1), nil),
		scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil),
		mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge),
		mocks.scalingTicketCount.EXPECT().WithLabelValues("applied").Return(scalingTicketCounter),
		mocks.scalingDurationSeconds.EXPECT().Observe(gomock.Any()),
		mocks.scaleResultCounter.EXPECT().WithLabelValues("ignored").Return(doneCounter),
	)
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 10}
	scaler, err := New(scaTgt, sObj, metrics, DryRunMode(true))
	require.NoError(t, err)
	require.NotNil(t, scaler)

	ticket := NewScalingTicket(5, false)
	scaler.applyScaleTicket(ticket)
	assert.False(t, scaler.desiredScale.isKnown)
	assert.Equal(t, uint(0), scaler.desiredScale.value)
}

func Test_ApplyScalingTicket_NoScaleObjectWatcherInDryRunMode(t *testing.T) {
	// This test was added to ensure that the ScaleObjectWatcher does not
	// run in dry-run mode. Why, see: https://github.com/ThomasObenaus/sokar/issues/98.

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, _ := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	sObj := ScalingObject{Name: "any", MinCount: 1, MaxCount: 10}
	scaler, err := New(scaTgt, sObj, metrics, DryRunMode(true), WatcherInterval(time.Millisecond*100))
	require.NoError(t, err)
	require.NotNil(t, scaler)

	scaler.Run()
	defer func() {
		err := scaler.Stop()
		assert.NoError(t, err)
		scaler.Join()
	}()

	// give the (potential) watcher some time
	time.Sleep(time.Millisecond * 200)

	// hint: This test would fail in case a running ScaleObjectWatcher would
	// call a method of the mocked ScalingTarget (e.g. GetCount or Scale)
}

func Test_ApplyScalingTicket_Scale(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	plannedButSkippedGauge := mock_metrics.NewMockGauge(mockCtrl)
	plannedButSkippedGauge.EXPECT().Set(float64(0))
	scalingTicketCounter := mock_metrics.NewMockCounter(mockCtrl)
	scalingTicketCounter.EXPECT().Inc()
	doneCounter := mock_metrics.NewMockCounter(mockCtrl)
	doneCounter.EXPECT().Inc()
	sObjName := "any"

	gomock.InOrder(
		scaTgt.EXPECT().GetScalingObjectCount(sObjName).Return(uint(0), nil),
		scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(false, nil),
		mocks.plannedButSkippedScalingOpen.EXPECT().WithLabelValues("UP").Return(plannedButSkippedGauge),
		scaTgt.EXPECT().AdjustScalingObjectCount(sObjName, uint(1), uint(10), uint(0), uint(5)).Return(nil),
		mocks.scalingTicketCount.EXPECT().WithLabelValues("applied").Return(scalingTicketCounter),
		mocks.scalingDurationSeconds.EXPECT().Observe(gomock.Any()),
		mocks.scaleResultCounter.EXPECT().WithLabelValues("done").Return(doneCounter),
	)
	sObj := ScalingObject{Name: sObjName, MinCount: 1, MaxCount: 10}
	scaler, err := New(scaTgt, sObj, metrics)
	require.NoError(t, err)
	require.NotNil(t, scaler)

	ticket := NewScalingTicket(5, false)
	scaler.applyScaleTicket(ticket)
	assert.True(t, scaler.desiredScale.isKnown)
	assert.Equal(t, uint(5), scaler.desiredScale.value)
}

func Test_OpenAndApplyScalingTicket(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	metrics, mocks := NewMockedMetrics(mockCtrl)

	sObjName := "any"
	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)
	scaTgt.EXPECT().GetScalingObjectCount(sObjName).Return(uint(0), nil).MaxTimes(11)
	scaTgt.EXPECT().IsScalingObjectDead(sObjName).Return(true, nil).MaxTimes(11)

	sObj := ScalingObject{Name: sObjName}
	scaler, err := New(scaTgt, sObj, metrics, MaxOpenScalingTickets(10))
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

func Test_UpdateDesiredScale(t *testing.T) {
	err := updateDesiredScale(scaleResult{}, nil)
	assert.Error(t, err)

	desiredScale := optionalValue{}
	err = updateDesiredScale(scaleResult{}, &desiredScale)
	assert.NoError(t, err)
	assert.False(t, desiredScale.isKnown)
	assert.Equal(t, uint(0), desiredScale.value)

	desiredScale = optionalValue{}
	desiredScale.setValue(10)
	err = updateDesiredScale(scaleResult{}, &desiredScale)
	assert.NoError(t, err)
	assert.True(t, desiredScale.isKnown)
	assert.Equal(t, uint(10), desiredScale.value)

	desiredScale = optionalValue{}
	err = updateDesiredScale(scaleResult{state: scaleDone, newCount: 10}, &desiredScale)
	assert.NoError(t, err)
	assert.True(t, desiredScale.isKnown)
	assert.Equal(t, uint(10), desiredScale.value)
}
