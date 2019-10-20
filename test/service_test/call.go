package main

import (
	"sync"

	gomock "github.com/golang/mock/gomock"
)

type Call interface {
	Times(n int) Call
	Return(rets ...interface{}) Call
	After(preReq Call) Call
}

type callImpl struct {
	gomockCall *gomock.Call
	wg         *sync.WaitGroup
}

func (c *callImpl) Times(n int) Call {
	c.gomockCall.Times(n)

	incWg := n - 1
	c.wg.Add(incWg)

	return c
}

func (c *callImpl) Return(rets ...interface{}) Call {
	c.gomockCall.Return(rets...)
	return c
}

func (c *callImpl) After(preReq Call) Call {

	preReqCall, _ := preReq.(*callImpl)
	// TODO: Check for failing cast

	c.gomockCall.After(preReqCall.gomockCall)
	return c
}

func InOrder(calls ...Call) {
	for i := 1; i < len(calls); i++ {
		calls[i].After(calls[i-1])
	}
}
