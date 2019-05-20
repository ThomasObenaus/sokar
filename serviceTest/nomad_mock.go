package serviceTest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type nomadMock struct {
	err error

	srv      *http.Server
	stopChan chan struct{}
}

func (nm *nomadMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Request %v\n", *r)

	w.WriteHeader(http.StatusOK)
}

func (nm *nomadMock) stop() {

	// context: wait for 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := nm.srv.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

func (nm *nomadMock) start() {
	nm.srv = &http.Server{
		Addr:    ":12000",
		Handler: nm,
	}

	// Run listening for messages in background
	go func() {
		err := nm.srv.ListenAndServe()

		if err != nil && err == http.ErrServerClosed {
			fmt.Println("API Srv shut down gracefully")
		} else {
			fmt.Printf("Failed serving: %s.\n", err.Error())
		}

		// send the stop message
		nm.stopChan <- struct{}{}
	}()
}
