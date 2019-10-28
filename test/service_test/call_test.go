package main

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_InOrder(t *testing.T) {

	InOrder(time.Now())

	mockCtrl := gomock.NewController(t)

	startTime := time.Now()
	defaultDeadline := time.Now()
	timeout1 := time.Second * 1
	call1 := &callImpl{timeout: timeout1, gomockCall: newGomockCall(mockCtrl), deadline: defaultDeadline}
	timeout2 := time.Second * 2
	call2 := &callImpl{timeout: timeout2, gomockCall: newGomockCall(mockCtrl), deadline: defaultDeadline}

	InOrder(
		startTime,
		call1,
		call2,
	)

	assert.Equal(t, startTime.Add(timeout1), call1.deadline)
	assert.Equal(t, defaultDeadline, call2.deadline)
}

func Test_UpdateDeadline(t *testing.T) {
	startTime := time.Now()
	timeout := time.Second * 10
	expectedDeadline := startTime.Add(timeout)

	call := callImpl{timeout: timeout}
	deadline := call.updateDeadline(startTime)
	assert.Equal(t, expectedDeadline, deadline)
}
