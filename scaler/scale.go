package scaler

import (
	"fmt"

	"github.com/thomasobenaus/sokar/helper"
)

// scaleState represents the state of a scaling
type scaleState string

const (
	// scaleUnknown means the scale process was completed successfully
	scaleUnknown scaleState = "unknown"
	// scaleDone means the scale process was completed successfully
	scaleDone scaleState = "done"
	// scaleRunning means the scale process is in progress
	scaleRunning scaleState = "running"
	// scaleFailed means the scale process was completed but failed
	scaleFailed scaleState = "failed"
	// scaleIgnored means the scale process was ignored (eventually not needed)
	scaleIgnored scaleState = "ignored"
	// scaleNotStarted means the scale process was not started yet
	scaleNotStarted scaleState = "not started"
)

// ScaleResult is created after scaling was done and contains the result
type scaleResult struct {
	state            scaleState
	stateDescription string
	newCount         uint
}

type policyCheckResult struct {
	validCount        uint
	desiredCount      uint
	minPolicyViolated bool
	maxPolicyViolated bool
}

func amountToScaleType(amount int) string {
	scaleTypeStr := "UP"
	if amount < 0 {
		scaleTypeStr = "DOWN"
	}
	return scaleTypeStr
}

// checkScalingPolicy verifies if the desired
func checkScalingPolicy(desiredCount uint, min uint, max uint) policyCheckResult {

	result := policyCheckResult{minPolicyViolated: false, maxPolicyViolated: false}

	result.desiredCount = desiredCount
	result.validCount = desiredCount

	// check if desiredCount exceeds minimum
	if desiredCount < min {
		result.validCount = min
		result.minPolicyViolated = true
	}

	// check if desiredCount exceeds maximum
	if desiredCount > max {
		result.validCount = max
		result.maxPolicyViolated = true
	}

	return result
}

// trueIfNil returns a scaleResult filled in with an appropriate error message in case the given scaler is nil
func trueIfNil(s *Scaler) (result scaleResult, ok bool) {
	ok = false
	result = scaleResult{state: scaleUnknown}

	if s == nil {
		ok = true
		result = scaleResult{
			state:            scaleFailed,
			stateDescription: "Scaler is nil",
			newCount:         0,
		}
	}
	return result, ok
}

// scale scales the job from currentCount to desiredCount.
// Internally it is checked if a scaling is needed and if the scaling policy is valid.
func (s *Scaler) scale(desiredCount uint, currentCount uint) scaleResult {
	if r, ok := trueIfNil(s); ok {
		return r
	}

	jobName := s.job.jobName
	min := s.job.minCount
	max := s.job.maxCount

	s.logger.Info().Str("job", jobName).Msgf("Request to scale job from %d to %d.", currentCount, desiredCount)

	dead, err := s.scalingTarget.IsJobDead(jobName)
	if err != nil {
		return scaleResult{
			state:            scaleFailed,
			stateDescription: fmt.Sprintf("Error obtaining if job is dead: %s.", err.Error()),
		}
	}

	if dead {
		return scaleResult{
			state:            scaleIgnored,
			stateDescription: fmt.Sprintf("Job '%s' is dead. Can't scale", jobName),
		}
	}

	chkResult := checkScalingPolicy(desiredCount, min, max)
	newCount := chkResult.validCount
	if chkResult.minPolicyViolated {
		s.logger.Info().Str("job", jobName).Msgf("Job.MinCount (%d) policy violated (wanted %d, have %d). Scale limited to %d.", min, chkResult.desiredCount, currentCount, newCount)
		s.metrics.scalingPolicyViolated.WithLabelValues("min").Inc()
	}
	if chkResult.maxPolicyViolated {
		s.logger.Info().Str("job", jobName).Msgf("Job.MaxCount (%d) policy violated (wanted %d, have %d). Scale limited to %d.", max, chkResult.desiredCount, currentCount, newCount)
		s.metrics.scalingPolicyViolated.WithLabelValues("max").Inc()
	}

	diff := helper.SubUint(newCount, currentCount)
	scaleNeeded := (diff != 0)

	if !scaleNeeded {
		s.logger.Info().Str("job", jobName).Msg("No scaling needed/ possible.")
		return scaleResult{
			state:            scaleIgnored,
			stateDescription: "No scaling needed/ possible.",
			newCount:         newCount,
		}
	}

	scaleTypeStr := amountToScaleType(diff)
	s.logger.Info().Str("job", jobName).Msgf("Scale %s by %d to %d.", scaleTypeStr, diff, newCount)

	// Set the new job count
	err = s.scalingTarget.SetJobCount(s.job.jobName, newCount)
	if err != nil {
		return scaleResult{
			state:            scaleFailed,
			stateDescription: fmt.Sprintf("Error adjusting job count to %d: %s.", newCount, err.Error()),
		}
	}

	return scaleResult{
		state:            scaleDone,
		stateDescription: "Scaling successfully done.",
		newCount:         newCount,
	}
}
