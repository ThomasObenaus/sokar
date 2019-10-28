package main

import (
	"fmt"
	reflect "reflect"
	"sync"
	"time"

	gomock "github.com/golang/mock/gomock"
)

// Call is a interface representing an expected call to the HTTP mock server
type Call interface {
	Return(rets ...interface{}) Call
	After(preReq Call) Call
	Within(timeout time.Duration) Call
	String() string

	// Internal methods
	join() (deadlineExpired bool)
	updateDeadline(start time.Time) time.Time
	release()
	commitCall()
}

type callImpl struct {
	// the underlying gomock call
	gomockCall *gomock.Call

	// the waitgroup needed to block/ wait until the expected call has arrived
	wg       *sync.WaitGroup
	timeout  time.Duration
	deadline time.Time

	// The succeeding call (can be nil)
	successor Call

	once sync.Once
}

func (c *callImpl) String() string {
	return c.gomockCall.String() + fmt.Sprintf(" (deadline=%s, timeout=%s)", c.deadline, c.timeout)
}

func (c *callImpl) release() {
	c.once.Do(func() {
		c.wg.Done()
	})
}

// join blocks until the expected call to the end point was made but at max until the internal
// deadline as been expired. In this case the method returns true.
func (c *callImpl) join() (deadlineExpired bool) {
	waitUntil(c.wg, c.deadline)

	return time.Now().After(c.deadline)
}

func NewCall(gomockCall *gomock.Call, method Method) Call {

	wg := &sync.WaitGroup{}
	wg.Add(1)

	call := callImpl{
		gomockCall: gomockCall,
		wg:         wg,
		timeout:    time.Second * 10,
		deadline:   time.Now().Add(time.Hour * 24), // initialized with a value that should never pass during test
		successor:  nil,
	}

	method(&call)

	return &call
}

// commitCall should be called as soon as the expected request has been triggered
func (c *callImpl) commitCall() {
	c.release()

	// Update the deadline of the succeeding call (if any)
	// as soon as this call was called.
	if c.successor != nil {
		c.successor.updateDeadline(time.Now())
	}
}

func (c *callImpl) Within(timeout time.Duration) Call {
	c.timeout = timeout
	return c
}

func (c *callImpl) Return(rets ...interface{}) Call {
	c.gomockCall.Return(rets...)
	return c
}

func (c *callImpl) After(preReq Call) Call {

	preReqCall, ok := preReq.(*callImpl)
	if !ok {
		panic(fmt.Sprintf("Failed to cast %s to %s", reflect.TypeOf(preReqCall).String(), reflect.TypeOf(&callImpl{}).String()))
	}

	c.gomockCall.After(preReqCall.gomockCall)

	// this callImpl is the successor of the given preReq
	preReqCall.successor = c
	return c
}

func (c *callImpl) updateDeadline(start time.Time) time.Time {
	c.deadline = start.Add(c.timeout)
	return c.deadline
}

func InOrder(startTime time.Time, calls ...Call) {

	if len(calls) == 0 {
		return
	}

	// Set the deadline for the first call in the order.
	// The deadline for the succeeding calls will be updated as soon
	// as the preceding call is called.
	firstCall := calls[0]
	firstCall.updateDeadline(startTime)

	// Store the order of the calls
	for i := 1; i < len(calls); i++ {
		predecessor := calls[i-1]
		call := calls[i]
		call.After(predecessor)
	}
}
