package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thomasobenaus/sokar/test/integration/helper"
	"github.com/thomasobenaus/sokar/test/integration/nomad"
)

func TestSimple(t *testing.T) {
	testCase := "Simple"
	sokarAddr := "http://localhost:11000"
	nomadAddr := "http://localhost:4646"
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
	//defer helper.SendScaleAlert("AlertA", false)
	NewJobAsserter(t, jobName, time.Millisecond*500, 30).AssertJobCount(34, jobName, d.GetJobCount)

}

type JobAsserter struct {
	t        *testing.T
	jobName  string
	waitTime time.Duration
	maxTries int
}

func NewJobAsserter(t *testing.T, jobName string, waitTime time.Duration, maxTries int) *JobAsserter {
	return &JobAsserter{
		t:        t,
		jobName:  jobName,
		waitTime: waitTime,
		maxTries: maxTries,
	}
}

func (ja *JobAsserter) AssertJobCount(expectedJobCount int, jobName string, obtainJobCountFunc func(jobName string) (int, error)) {
	count := 0
	for i := 0; i < ja.maxTries; i++ {
		var err error
		count, err = obtainJobCountFunc(ja.jobName)
		require.NoError(ja.t, err, "Failed to obtain job count")
		if count == expectedJobCount {
			return // success case -> no assert, just return
		}
		ja.t.Logf("Jobcount not as expected (%d), but was %d. Recheck in %s\n", expectedJobCount, count, ja.waitTime)
		time.Sleep(ja.waitTime)
	}

	assert.Failf(ja.t, "Jobcount invalid", "Jobcount is not %d as expected but was %d at last try (#tries=%d).", expectedJobCount, count, ja.maxTries)
}
