package main

import (
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
}
