package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/thomasobenaus/sokar/api"
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

	ctrl := NewController(t)
	defer ctrl.Finish()

	nomadMock := NewMockHTTP(ctrl, 18000, FailOnUnexpectedCalls(true))
	awsMock := NewMockHTTP(ctrl, 18001, FailOnUnexpectedCalls(true))

	InOrder(
		time.Now(),
		nomadMock.EXPECT().GET("/health").Within(time.Second*10).Return(NewStringResponse("{\"Call\": 1}", AddHeader("Content-Type", "application/json"))),
		awsMock.EXPECT().GET("/health").Within(time.Second*10).Return(NewStringResponse("{\"Call\": 2}", AddHeader("Content-Type", "application/protobuf"))),
		nomadMock.EXPECT().GET("/health").Within(time.Second*10).Return(NewStringResponse("{\"Call\": 3}", AddHeader("Content-Type", "application/json"))),
	)
}

func Test_Job(t *testing.T) {

	ctrl := NewController(t)
	defer ctrl.Finish()

	nomadMock := NewMockHTTP(ctrl, 18000, FailOnUnexpectedCalls(false))

	InOrder(
		time.Now(),
		nomadMock.EXPECT().GET("/v1/job/fail-service").Within(time.Second*50).Return(NewStringResponse("{\"Call\": 1}", AddHeader("Content-Type", "application/json"))),
	)

	//sendAlert("http://127.0.0.1:11000/api/alerts")
}

func Test_Jow(t *testing.T) {

	server := api.New(18000)
	server.Run()

	counter := 0

	server.Router.HandlerFunc("GET", "/v1/job/fail-service", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%d] Request %v\n", counter, r)
		counter++

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"Call\": 1}"))
	})

	time.Sleep(time.Second * 30)
}

func sendAlert(url string) {
	data := []byte(`{
		"receiver": "PM",
		"status": "firing",
		"alerts": [
		  {
			"status": "firing",
			"labels": {
			  "alertname": "AlertA",
			  "alert-type": "scaling",
			  "scale-type": "up"
			},
			"annotations": {
			  "description": "Scales the component XYZ UP"
			},
			"startsAt": "2019-02-23T12:00:00.000+01:00",
			"endsAt": "2019-02-23T12:05:00.000+01:00",
			"generatorURL": "http://generator_url"
		  },
		  {
			"status": "firings",
			"labels": {
			  "alertname": "AlertB",
			  "alert-type": "scaling",
			  "scale-type": "down"
			},
			"annotations": {
			  "description": "Scales the component XYZ DOWN"
			},
			"startsAt": "2019-02-23T12:00:00.000+01:00",
			"endsAt": "2019-02-23T12:05:00.000+01:00",
			"generatorURL": "http://generatorURL"
		  }
		],
		"groupLabels": {},
		"commonLabels": { "alertname": "AlertA" },
		"commonAnnotations": {},
		"externalURL": "http://externalURL",
		"version": "4",
		"groupKey": "{}:{}"
	  }`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}
