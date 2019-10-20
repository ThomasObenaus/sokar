package main

import (
	"fmt"
	"sync"
	"time"
)

func wait(wg *sync.WaitGroup, timeout time.Duration) {

	deadline := time.Now().Add(timeout)
	waitFor := make(chan struct{}, 0)

	go func() {
		defer close(waitFor)
		wg.Wait()
	}()

	select {
	case <-time.After(timeout):
		fmt.Printf("--[Dbg] Time Expired (timeout=%s)\n", timeout.String())
		return
	case <-waitFor:
		fmt.Printf("--[Dbg] WaitGroup released (time left=%s)\n", deadline.Sub(time.Now()).String())
		return
	}
}
