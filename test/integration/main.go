package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	sokarAddr := "http://localhost:11000"
	nomadAddr := "http://localhost:4646"

	fmt.Printf("Start waiting for nomad (%s)....\n", nomadAddr)
	internalIP, err := waitForNomad(nomadAddr, time.Second*2, 20)
	if err != nil {
		log.Fatalf("Failed while waiting for nomad: %s\n", err.Error())
	}
	fmt.Printf("Nomad up and running (internal-ip=%s)\n", internalIP)

	fmt.Printf("Start waiting for sokar (%s)....\n", sokarAddr)
	err = waitForSokar(sokarAddr, time.Second*2, 20)
	if err != nil {
		log.Fatalf("Failed while waiting for sokar: %s\n", err.Error())
	}
	fmt.Printf("Sokar up and running\n")

	fmt.Printf("Deploy Job\n")
	d, err := NewDeployer(nomadAddr)
	if err != nil {
		log.Fatalf("Failed to create deployer: %s\n", err.Error())
	}

	job := NewJobDescription("fail-service", "testing", "thobe/fail_service:v0.1.0", 2)
	err = d.Deploy(job)
	if err != nil {
		log.Fatalf("Failed to deploy job: %s\n", err.Error())
	}

	fmt.Printf("Deploy Job succeeded\n")
}

func waitForNomad(nomadAddr string, timeoutBetweenTries time.Duration, numTries int) (string, error) {
	queryPath := "/v1/status/leader"
	logPrefix := "wait for nomad"
	return waitForService(nomadAddr, queryPath, logPrefix, timeoutBetweenTries, numTries)
}

func waitForSokar(serviceAddr string, timeoutBetweenTries time.Duration, numTries int) error {
	queryPath := "/health"
	logPrefix := "wait for sokar"
	_, err := waitForService(serviceAddr, queryPath, logPrefix, timeoutBetweenTries, numTries)
	return err
}

func waitForService(serviceAddr, queryPath, logPrefix string, timeoutBetweenTries time.Duration, numTries int) (string, error) {
	client := http.Client{
		Timeout: time.Millisecond * 500,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	queryURL := fmt.Sprintf("%s%s", serviceAddr, queryPath)
	for i := 0; i < numTries; i++ {
		fmt.Printf("[%s] %d/%d\n", logPrefix, i+1, numTries)
		if i > 0 {
			time.Sleep(timeoutBetweenTries)
		}
		resp, err := client.Get(queryURL)
		if err != nil {
			fmt.Printf("[%s] %s\n", logPrefix, err.Error())
			continue
		}

		if resp == nil {
			fmt.Printf("[%s] Response is nil\n", logPrefix)
			continue
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("[%s] Error reading response %s\n", logPrefix, err.Error())
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("[%s] Error querying service at '%s' got response [%d] '%s'", logPrefix, serviceAddr, resp.StatusCode, string(bodyBytes))
			continue
		}

		return string(bodyBytes), nil
	}

	return "", fmt.Errorf("[%s] service not running at %s (timeoutBetweenTries=%s, numTries=%d)", logPrefix, serviceAddr, timeoutBetweenTries, numTries)
}
