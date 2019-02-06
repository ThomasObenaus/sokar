package scaler

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
}

func TestScaleBy_JobDead(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname}
	scaler, err := cfg.New(scaTgt)

	// dead job - error
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, fmt.Errorf("internal error"))
	err = scaler.ScaleBy(2)
	assert.Error(t, err)

	// dead job
	scaTgt.EXPECT().IsJobDead(jobname).Return(true, nil)
	err = scaler.ScaleBy(2)
	assert.Error(t, err)
}

func TestScaleBy_Up(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)

	// scale up
	currentJobCount := uint(0)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(2)).Return(nil)
	err = scaler.ScaleBy(2)
	assert.NoError(t, err)

	// scale up - relative
	currentJobCount = uint(1)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(3)).Return(nil)
	err = scaler.ScaleBy(2)
	assert.NoError(t, err)

	// scale up - max hit
	currentJobCount = uint(4)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(5)).Return(nil)
	err = scaler.ScaleBy(2)
	assert.NoError(t, err)
}

func TestScaleBy_Down(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)

	// scale down
	currentJobCount := uint(3)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(1)).Return(nil)
	err = scaler.ScaleBy(-2)
	assert.NoError(t, err)

	// scale up - min hit
	currentJobCount = uint(2)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	scaTgt.EXPECT().SetJobCount(jobname, uint(1)).Return(nil)
	err = scaler.ScaleBy(-5)
	assert.NoError(t, err)
}

func TestScaleBy_NoScale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scaTgt := mock_scaler.NewMockScalingTarget(mockCtrl)

	jobname := "any"
	cfg := Config{JobName: jobname, MinCount: 1, MaxCount: 5}
	scaler, err := cfg.New(scaTgt)

	// scale down
	currentJobCount := uint(5)
	scaTgt.EXPECT().IsJobDead(jobname).Return(false, nil)
	scaTgt.EXPECT().GetJobCount(jobname).Return(currentJobCount, nil)
	err = scaler.ScaleBy(2)
	assert.NoError(t, err)
}
