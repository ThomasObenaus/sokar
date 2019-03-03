package scaler

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sokar "github.com/thomasobenaus/sokar/sokar/iface"
	"github.com/thomasobenaus/sokar/test/scaler"
)

func TestNew(t *testing.T) {

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
	assert.Nil(t, scaler.scalingTicket)
	assert.NotNil(t, scaler.scalingTarget)
}

func TestScaleTo(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// scale up
	currentJobCount := uint(0)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(2)).Return(nil)
	result := scaler.ScaleTo(2)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)

	// scale err
	scaTgt.EXPECT().GetJobCount(jobname).Return(uint(0), fmt.Errorf("internal err"))
	result = scaler.ScaleTo(2)
	assert.Equal(t, sokar.ScaleFailed, result.State)
}

func TestScaleBy(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)
	require.NoError(t, err)

	// scale up
	currentJobCount := uint(0)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(2)).Return(nil)
	result := scaler.ScaleBy(2)
	assert.NotEqual(t, sokar.ScaleFailed, result.State)

	// scale err
	scaTgt.EXPECT().GetJobCount(jobname).Return(uint(0), fmt.Errorf("internal err"))
	result = scaler.ScaleBy(2)
	assert.Equal(t, sokar.ScaleFailed, result.State)
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
