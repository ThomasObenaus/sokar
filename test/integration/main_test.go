package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/thomasobenaus/sokar/test/integration/helper"
	"github.com/thomasobenaus/sokar/test/integration/nomad"
)

func TestScaleUp(t *testing.T) {
	testCase := "ScaleUp"
	sokarAddr := helper.SokarAddr
	nomadAddr := helper.NomadAddr
	jobName := "fail-service"

	helper.PrintCheckPoint(testCase, "Start waiting for nomad (%s)....\n", nomadAddr)
	internalIP, err := helper.WaitForNomad(t, nomadAddr, time.Second*2, 20)
	require.NoError(t, err, "Failed while waiting for nomad")

	helper.PrintCheckPoint(testCase, "Nomad up and running (internal-ip=%s)\n", internalIP)

	helper.PrintCheckPoint(testCase, "Start waiting for sokar (%s)....\n", sokarAddr)
	err = helper.WaitForSokar(t, sokarAddr, time.Second*2, 20)
	require.NoError(t, err, "Failed while waiting for sokar")
	helper.PrintCheckPoint(testCase, "Sokar up and running\n")

	helper.PrintCheckPoint(testCase, "Deploy Job\n")
	d, err := nomad.NewDeployer(t, nomadAddr)
	require.NoError(t, err, "Failed to create deployer")

	job := nomad.NewJobDescription(jobName, "testing", "thobe/fail_service:v0.1.0", 2, map[string]string{"HEALTHY_FOR": "-1"})
	err = d.Deploy(job)
	require.NoError(t, err, "Failed to deploy job")

	count, err := d.GetJobCount(jobName)
	require.NoError(t, err, "Failed to obtain job count")
	require.Equal(t, 2, count, "Job count not as expected after initial deployment")

	helper.PrintCheckPoint(testCase, "Deploy Job succeeded\n")

	helper.PrintCheckPoint(testCase, "Sending scale alert\n")
	err = helper.SendScaleAlert("AlertA", true)
	require.NoError(t, err)
	helper.PrintCheckPoint(testCase, "Scale alert sent\n")

	// ensure to disable the alert from firing
	defer helper.SendScaleAlert("AlertA", false)

	helper.PrintCheckPoint(testCase, "Check if job was scaled to the expected count\n")
	helper.NewJobAsserter(t, jobName, time.Millisecond*500, 50).AssertJobCount(3, jobName, d.GetJobCount)
}
