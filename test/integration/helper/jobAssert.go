package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// JobAsserter helps on different assertion tasks. Especially those where multiple retries have to be done before
// considering a test as failed.
type JobAsserter struct {
	t        *testing.T
	jobName  string
	waitTime time.Duration // time to wait between the tries
	maxTries int           // max number of tries to get the assertion true
}

// NewJobAsserter creates a new JobAsserter
func NewJobAsserter(t *testing.T, jobName string, waitTime time.Duration, maxTries int) *JobAsserter {
	return &JobAsserter{
		t:        t,
		jobName:  jobName,
		waitTime: waitTime,
		maxTries: maxTries,
	}
}

// AssertJobCount asserts that the count of the nomad job is as expected
func (ja *JobAsserter) AssertJobCount(expectedJobCount int, jobName string, obtainJobCountFunc func(jobName string) (int, error)) {
	count := 0
	for i := 0; i < ja.maxTries; i++ {
		var err error
		count, err = obtainJobCountFunc(ja.jobName)
		require.NoError(ja.t, err, "Failed to obtain job count")
		if count == expectedJobCount {
			ja.t.Logf("Jobcount as expected: %d==%d\n", expectedJobCount, count)
			return // success case -> no assert, just return
		}
		ja.t.Logf("Jobcount not as expected (%d), but was %d. Recheck in %s\n", expectedJobCount, count, ja.waitTime)
		time.Sleep(ja.waitTime)
	}

	assert.Failf(ja.t, "Jobcount invalid", "Jobcount is not %d as expected but was %d at last try (#tries=%d).", expectedJobCount, count, ja.maxTries)
}
