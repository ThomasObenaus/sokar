package main

import (
	"context"
	"fmt"
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

	//ctx := myCTX{donechan: make(chan struct{}, 1)}

	//mockCtrl := gomock.NewController(t)

	ctx := context.Background()
	//ctx, cn := context.WithTimeout(ctx, time.Second*10)
	//_ = cn

	mockCtrl, ctx := gomock.WithContext(ctx, t)

	go func() {
		fmt.Println("STAAAAAAAART")

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("EEEEEEEEEEEEEEEND %s\n", ctx.Err().Error())
				return
			}
		}
	}()

	//mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := NewMockHTTP(mockCtrl, 18000)
	_ = mock
	//mock.GET("jj")

	mock.EXPECT().GET("/health").Return(http.StatusOK, "BLA")
	//mock.EXPECT().GET("/api/alerts").Return(http.StatusNotFound, "BLA")

	//server := New(18000, t)
	//defer server.Close()
	//
	//server.GET("/health", http.StatusOK, "SUPER")
	//
	time.Sleep(time.Second * 5)

	//mockCtrl := gomock.NewController(t)
	//defer mockCtrl.Finish()
	//mock := NewMockHTTP(mockCtrl)
	//
	//mock.EXPECT().POST("HELLO").Return(http.StatusOK, "huhu")
	//
	//receiver := api.New(18000)
	//receiver.Run()
	//// receiver.GET(/health,response)
	//receiver.Router.HandlerFunc("GET", "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//
	//	code, _ := mock.POST("HELLO")
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(code)
	//	io.WriteString(w, "data")
	//
	//	//receiver.Stop()
	//}))
	//time.Sleep(time.Millisecond * 100)
	//
	//res, err := http.Get("http://127.0.0.1:18000/health")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//greeting, err := ioutil.ReadAll(res.Body)
	//res.Body.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Printf("%s\n", greeting)
	//
	////receiver.Join()
	//receiver.Stop()

}
