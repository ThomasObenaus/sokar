package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

// # Notes:
// - Durations timeouts are optional, since the test itself has a defined deadline
// ## Procedure
// Sokar is at http://127.0.0.1:11000"
// 1. Send Request
// http.POST("http://127.0.0.1:11000/api/alerts",JSON{Alert{firing}})
// 2. Wait for Request to nomad, expect a certain body and respond with suitable data
// expect(timeout time.Duration).POST(data JSON).Return(code http.StatusCode, data JSON)

func Test_My(t *testing.T) {

	mock := WithTimeout(t, 18000, time.Second*180)
	defer mock.Finish()

	//mock.EXPECT().GET("/health").Return(http.StatusOK, "BLA")
	gomock.InOrder(
		mock.EXPECT().GET("/v1/job/fail-service").Return(http.StatusOK, "BLA"),
	)

	// TODO: Solve Times(n) problem
	// Probably the mock methods should return an interface aligned to gomock.Call
	// and then wrap the Times(n) method by internally incrementing/ setting the wg accordingly
	// e.g.
	// func (c *myCall) Times(n int) *Call{
	// c.wg.Add(n)
	// return c.gomockCall
	//}
	//
	// Works probably also for the in order problem
}
