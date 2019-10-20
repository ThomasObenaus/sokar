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

	mock := WithTimeout(t, 18000, time.Second*30)
	defer mock.Finish()

	//mock.EXPECT().GET("/health").Return(http.StatusOK, "BLA")
	gomock.InOrder(
		mock.EXPECT().GET("/health").Return(http.StatusOK, "BLA"),
		mock.EXPECT().GET("/health").Return(http.StatusOK, "BLUBB"),
		mock.EXPECT().GET("/healths").Return(http.StatusOK, "BLUBB2"),
	)

}
