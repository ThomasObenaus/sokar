package main

import (
	"net/http"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
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

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := NewMockHTTP(mockCtrl, 18000)
	mock.EXPECT().GET("/health").Return(http.StatusOK, "BLA")

	time.Sleep(time.Second * 5)
}
