package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/thomasobenaus/sokar/test/integration/helper"
	"github.com/thomasobenaus/sokar/test/integration/nomad"
)

func TestSimple(t *testing.T) {
	sokarAddr := "http://localhost:11000"
	nomadAddr := "http://localhost:4646"
	jobName := "fail-service"

	t.Logf("Start waiting for nomad (%s)....\n", nomadAddr)
	internalIP, err := helper.WaitForNomad(t, nomadAddr, time.Second*2, 20)
	require.NoError(t, err, "Failed while waiting for nomad")

	t.Logf("Nomad up and running (internal-ip=%s)\n", internalIP)

	t.Logf("Start waiting for sokar (%s)....\n", sokarAddr)
	err = helper.WaitForSokar(t, sokarAddr, time.Second*2, 20)
	require.NoError(t, err, "Failed while waiting for sokar")
	t.Logf("Sokar up and running\n")

	t.Logf("Deploy Job\n")
	d, err := nomad.NewDeployer(t, nomadAddr)
	require.NoError(t, err, "Failed to create deployer")

	job := nomad.NewJobDescription(jobName, "testing", "thobe/fail_service:v0.1.0", 2, map[string]string{"HEALTHY_FOR": "-1"})
	err = d.Deploy(job)
	require.NoError(t, err, "Failed to deploy job")

	count, err := d.GetJobCount(jobName)
	require.NoError(t, err, "Failed to obtain job count")
	require.Equal(t, 2, count, "Job count not as expected after initial deployment")

	t.Logf("Deploy Job succeeded\n")

	sendScaleAlert()
}

func sendScaleAlert() {
	client := http.Client{
		Timeout: time.Millisecond * 500,
	}

	bodybytes := []byte(`{
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
	body := bytes.NewReader(bodybytes)
	resp, err := client.Post("http://localhost:11000/api/alerts", "application/json", body)

	if err != nil {
		log.Fatalf("Error sending request: %s\n", err.Error())
	}

	fmt.Printf("Request send, response: %v\n", resp)
}
