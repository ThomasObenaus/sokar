package main

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
)

type cleanupCallback func() error

type Controller struct {
	gmckCtrl         *gomock.Controller
	calls            []Call
	cleanupCallbacks []cleanupCallback
}

func NewController(t *testing.T) *Controller {

	ctrl := Controller{
		gmckCtrl:         gomock.NewController(t),
		calls:            make([]Call, 0),
		cleanupCallbacks: make([]cleanupCallback, 0),
	}

	return &ctrl
}

func (c *Controller) CreateAndRegisterGETCall(gomockCall *gomock.Call) Call {
	call := NewCall(gomockCall, GET())
	c.calls = append(c.calls, call)
	return call
}

// Finish has to be called at the end to clean up and to check if all expected calls where made.
func (c *Controller) Finish() {

	// Wait here for all registered calls until they succeed.
	// And fail immediately in case their deadline (timeout) has been exceeded.
	for _, call := range c.calls {
		deadlineIsExpired := call.join()
		if deadlineIsExpired {
			c.gmckCtrl.T.Fatalf("The deadline for call '%v' has been expired before someone called the according end-point.", call)
		}
	}

	// wait to give the latest responses some time to be read from the receiver
	time.Sleep(time.Millisecond * 200)

	// clean up
	c.gmckCtrl.Finish()

	// call all registered cleanupCallbacks
	for _, cleanupCb := range c.cleanupCallbacks {
		err := cleanupCb()

		if err != nil {
			c.gmckCtrl.T.Errorf("Error calling cleanupCallback: %s", err.Error())
		}
	}
}

func (c *Controller) releaseAllCallLocks() {
	for _, call := range c.calls {
		call.release()
	}
}

func (c *Controller) addCleanupCallback(cb cleanupCallback) {
	c.cleanupCallbacks = append(c.cleanupCallbacks, cb)
}
