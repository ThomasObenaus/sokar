package scaler

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/sokar/iface"
	"github.com/thomasobenaus/sokar/test/scaler"
)

func TestScale_JobDead(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// dead job - error
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, fmt.Errorf("internal error"))
	result := scaler.scale(2, 0)
	assert.Equal(t, sokar.ScaleFailed, result.State)

	// dead job
	scaTgt.EXPECT().IsJobDead(jobname).Return(true, nil)
	result = scaler.scale(2, 0)
	assert.Equal(t, sokar.ScaleIgnored, result.State)
}

func TestScale_Up(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// scale up
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(2)).Return(nil)
	result := scaler.scale(2, 0)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)

	// scale up - max hit
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(5)).Return(nil)
	result = scaler.scale(6, 0)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)
}

func TestScale_Down(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// scale down
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(1)).Return(nil)
	result := scaler.scale(1, 4)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)

	// scale up - min hit
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(1)).Return(nil)
	result = scaler.scale(0, 2)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)
}

func TestScale_NoScale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// scale down
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	result := scaler.scale(2, 2)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)
}

func TestScaleBy_CheckScalingPolicy(t *testing.T) {

	chk := checkScalingPolicy(0, 0, 0)
	assert.Equal(t, uint(0), chk.validCount)
	assert.Equal(t, uint(0), chk.desiredCount)
	assert.False(t, chk.minPolicyViolated)
	assert.False(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(1, 0, 0)
	assert.Equal(t, uint(0), chk.validCount)
	assert.Equal(t, uint(1), chk.desiredCount)
	assert.False(t, chk.minPolicyViolated)
	assert.True(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(1, 2, 3)
	assert.Equal(t, uint(2), chk.validCount)
	assert.Equal(t, uint(1), chk.desiredCount)
	assert.True(t, chk.minPolicyViolated)
	assert.False(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(3, 1, 4)
	assert.Equal(t, uint(3), chk.validCount)
	assert.Equal(t, uint(3), chk.desiredCount)
	assert.False(t, chk.minPolicyViolated)
	assert.False(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(3, 1, 2)
	assert.Equal(t, uint(2), chk.validCount)
	assert.Equal(t, uint(3), chk.desiredCount)
	assert.False(t, chk.minPolicyViolated)
	assert.True(t, chk.maxPolicyViolated)

	chk = checkScalingPolicy(2, 3, 1)
	assert.Equal(t, uint(1), chk.validCount)
	assert.Equal(t, uint(2), chk.desiredCount)
	assert.True(t, chk.minPolicyViolated)
	assert.True(t, chk.maxPolicyViolated)
}

func TestScaleBy_trueIfNil(t *testing.T) {
	_, ok := trueIfNil(nil)
	assert.True(t, ok)

	scaler := &Scaler{}
	_, ok = trueIfNil(scaler)
	assert.False(t, ok)
}