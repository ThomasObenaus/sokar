package main

import (
	"fmt"
	"sync"
	"time"
)

func Wait(wg *sync.WaitGroup, timeout time.Duration) {

	deadline := time.Now().Add(timeout)
	waitFor := make(chan struct{}, 0)

	go func() {
		defer close(waitFor)
		wg.Wait()
	}()

	select {
	case <-time.After(timeout):
		fmt.Printf("Time Expired (timeout=%s)\n", timeout.String())
		return
	case <-waitFor:
		fmt.Printf("Wg released (time left=%s)\n", deadline.Sub(time.Now()).String())
		return
	}
}

func DoneIn(wg *sync.WaitGroup, timeout time.Duration) {

	deadline := time.Now().Add(timeout)
	waitFor := make(chan struct{}, 0)

	go func() {
		defer close(waitFor)
		wg.Wait()
	}()

	go func() {
		select {
		case <-time.After(timeout):
			fmt.Printf("Time Expired (timeout=%s)\n", timeout.String())
			wg.Done()
			return
		case <-waitFor:
			fmt.Printf("Wg released (time left=%s)\n", deadline.Sub(time.Now()).String())
			return
		}
	}()
}
