package main

import (
	"net/http"
	"testing"
	"time"
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

	mock := NewMockHTTP(t, 18000, WithTimeout(time.Second*180), FailOnUnexpectedCalls(true))
	defer mock.Finish()

	//mock.EXPECT().GET("/health").Return(http.StatusOK, "BLA")
	InOrder(
		mock.EXPECT().GET("/health").Return(http.StatusOK, "BLA").Times(2),
		mock.EXPECT().GET("/healths").Return(http.StatusOK, "BLA").Times(1),
		mock.EXPECT().GET("/health").Return(http.StatusOK, "BLA").Times(3),
	)
}
