package helper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// WaitForNomad wait until nomad is up and running
func WaitForNomad(t *testing.T, nomadAddr string, timeoutBetweenTries time.Duration, numTries int) (string, error) {
	queryPath := "/v1/status/leader"
	logPrefix := "wait for nomad"
	return waitForService(t, nomadAddr, queryPath, logPrefix, timeoutBetweenTries, numTries)
}

// WaitForSokar wait until sokar is up and running
func WaitForSokar(t *testing.T, serviceAddr string, timeoutBetweenTries time.Duration, numTries int) error {
	queryPath := "/health"
	logPrefix := "wait for sokar"
	_, err := waitForService(t, serviceAddr, queryPath, logPrefix, timeoutBetweenTries, numTries)
	return err
}

func waitForService(t *testing.T, serviceAddr, queryPath, logPrefix string, timeoutBetweenTries time.Duration, numTries int) (string, error) {
	client := http.Client{
		Timeout: time.Millisecond * 500,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	queryURL := fmt.Sprintf("%s%s", serviceAddr, queryPath)
	for i := 0; i < numTries; i++ {
		t.Logf("[%s] %d/%d\n", logPrefix, i+1, numTries)
		if i > 0 {
			time.Sleep(timeoutBetweenTries)
		}
		resp, err := client.Get(queryURL)
		if err != nil {
			t.Logf("[%s] %s\n", logPrefix, err.Error())
			continue
		}

		if resp == nil {
			t.Logf("[%s] Response is nil\n", logPrefix)
			continue
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Logf("[%s] Error reading response %s\n", logPrefix, err.Error())
			continue
		}

		if resp.StatusCode != http.StatusOK {
			t.Logf("[%s] Error querying service at '%s' got response [%d] '%s'", logPrefix, serviceAddr, resp.StatusCode, string(bodyBytes))
			continue
		}

		return string(bodyBytes), nil
	}

	return "", fmt.Errorf("[%s] service not running at %s (timeoutBetweenTries=%s, numTries=%d)", logPrefix, serviceAddr, timeoutBetweenTries, numTries)
}
