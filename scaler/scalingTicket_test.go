package scaler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sokar "github.com/thomasobenaus/sokar/sokar/iface"
)

func TestNewScalingTicket(t *testing.T) {
	sj := NewScalingTicket(0)
	assert.WithinDuration(t, time.Now(), sj.issuedAt, time.Millisecond*100)
	assert.Equal(t, sokar.ScaleNotStarted, sj.state)
	assert.Equal(t, uint(0), sj.desiredCount)
	assert.Nil(t, sj.startedAt)
	assert.Nil(t, sj.completedAt)
}

func Test_Start(t *testing.T) {
	sj := NewScalingTicket(10)
	sj.start()
	require.NotNil(t, sj.startedAt)
	assert.WithinDuration(t, time.Now(), *sj.startedAt, time.Millisecond*100)
	assert.Equal(t, sokar.ScaleRunning, sj.state)
	assert.Nil(t, sj.completedAt)
}

func Test_Complete(t *testing.T) {
	sj := NewScalingTicket(10)
	sj.complete(sokar.ScaleDone)
	require.NotNil(t, sj.completedAt)
	assert.WithinDuration(t, time.Now(), *sj.completedAt, time.Millisecond*100)
	assert.Equal(t, sokar.ScaleDone, sj.state)
}

func Test_TicketAge(t *testing.T) {
	sj := NewScalingTicket(1)
	sj.issuedAt = time.Now().Add(time.Second * -10)
	assert.InDelta(t, time.Second*10, sj.ticketAge(), float64(time.Microsecond*100))
}

func Test_LeadDuration(t *testing.T) {
	sj := NewScalingTicket(1)

	_, err := sj.leadDuration()
	assert.Error(t, err)

	sj.issuedAt = time.Now().Add(time.Second * -10)
	in10Sec := time.Now().Add(time.Second * 10)
	sj.completedAt = &in10Sec
	dur, err := sj.leadDuration()
	assert.NoError(t, err)
	assert.InDelta(t, time.Second*20, dur, float64(time.Microsecond*100))
}

func Test_ProcessingDuration(t *testing.T) {
	sj := NewScalingTicket(1)

	_, err := sj.processingDuration()
	assert.Error(t, err)

	sj.start()
	_, err = sj.processingDuration()
	assert.Error(t, err)

	ago10Sec := time.Now().Add(time.Second * -10)
	sj.startedAt = &ago10Sec
	in10Sec := time.Now().Add(time.Second * 10)
	sj.completedAt = &in10Sec
	dur, err := sj.processingDuration()
	assert.NoError(t, err)
	assert.InDelta(t, time.Second*20, dur, float64(time.Microsecond*100))
}
