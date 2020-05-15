package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	nomadAddr := "http://localhost:4646"

	fmt.Printf("Start waiting for nomad (%s)....\n", nomadAddr)
	internalIP, err := waitForNomad(nomadAddr, time.Second*2, 20)
	if err != nil {
		log.Fatalf("Failed while waiting for nomad: %s\n", err.Error())
	}

	fmt.Printf("Nomad up and running (internal-ip=%s)\n", internalIP)
}

func waitForNomad(nomadAddr string, timeoutBetweenTries time.Duration, numTries int) (string, error) {

	client := http.Client{
		Timeout: time.Millisecond * 500,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	statusURL := fmt.Sprintf("%s/v1/status/leader", nomadAddr)
	for i := 0; i < numTries; i++ {
		fmt.Printf("[wait for nomad] %d/%d\n", i+1, numTries)
		if i > 0 {
			time.Sleep(timeoutBetweenTries)
		}
		resp, err := client.Get(statusURL)
		if err != nil {
			fmt.Printf("[wait for nomad] %s\n", err.Error())
			continue
		}

		if resp == nil {
			fmt.Println("[wait for nomad] Response is nil")
			continue
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("[wait for nomad] Error reading response %s\n", err.Error())
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("[wait for nomad] Error querying nomad at '%s' got response [%d] '%s'", nomadAddr, resp.StatusCode, string(bodyBytes))
			continue
		}

		return string(bodyBytes), nil
	}

	return "", fmt.Errorf("nomad not running at %s (timeoutBetweenTries=%s, numTries=%d)", nomadAddr, timeoutBetweenTries, numTries)
}
