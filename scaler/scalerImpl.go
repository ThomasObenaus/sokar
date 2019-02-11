package scaler

import (
	"fmt"
	"math"

	"github.com/thomasobenaus/sokar/sokar"
)

type jobConfig struct {
	jobName  string
	minCount uint
	maxCount uint
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

func checkScalingPolicyNew(count uint, min uint, max uint) policyCheckResult {

	result := policyCheckResult{minPolicyViolated: false, maxPolicyViolated: false}

	result.desiredCount = count

	// check if count exceeds minimum
	if count < min {
		count = min
		result.minPolicyViolated = true
	}

	// check if count exceeds maximum
	if count > max {
		count = max
		result.maxPolicyViolated = true
	}

	result.validCount = count

	return result
}

func checkScalingPolicy(count uint, amount int, min uint, max uint) policyCheckResult {

	result := policyCheckResult{minPolicyViolated: false, maxPolicyViolated: false}

	newCountTmp := int(count) + amount
	newCount := uint(newCountTmp)
	if newCountTmp < 0 {
		newCount = 0
	}

	result.desiredCount = newCount

	// check if count exceeds minimum
	if newCount < min {
		newCount = min
		result.minPolicyViolated = true
	}

	// check if count exceeds maximum
	if newCount > max {
		newCount = max
		result.maxPolicyViolated = true
	}

	result.validCount = newCount

	return result
}

// trueIfNil returns a ScaleResult filled in with an appropriate error message in case the given scaler is nil
func trueIfNil(s *Scaler) (result sokar.ScaleResult, ok bool) {
	ok = false
	result = sokar.ScaleResult{State: sokar.ScaleUnknown}

	if s == nil {
		ok = true
		result = sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: "Scaler is nil",
			NewCount:         0,
		}
	}
	return result, ok
}

func (s *Scaler) ScaleTo(count uint) sokar.ScaleResult {
	if r, ok := trueIfNil(s); ok {
		return r
	}

	jobName := s.job.jobName
	min := s.job.minCount
	max := s.job.maxCount

	s.logger.Info().Str("job", jobName).Msgf("Request to scale job to %d.", count)

	dead, err := s.scalingTarget.IsJobDead(jobName)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error obtaining if job is dead: %s.", err.Error()),
		}
	}

	if dead {
		return sokar.ScaleResult{
			State:            sokar.ScaleIgnored,
			StateDescription: fmt.Sprintf("Job '%s' is dead. Can't scale", jobName),
		}
	}

	chkResult := checkScalingPolicyNew(count, min, max)
	newCount := chkResult.validCount
	if chkResult.minPolicyViolated {
		s.logger.Info().Str("job", jobName).Msgf("Job.MinCount (%d) policy violated (wanted %d). Scale limited to %d.", min, chkResult.desiredCount, count, newCount)
	}
	if chkResult.maxPolicyViolated {
		s.logger.Info().Str("job", jobName).Msgf("Job.MinCount (%d) policy violated (wanted %d). Scale limited to %d.", min, chkResult.desiredCount, count, newCount)
	}

	diff := int(math.Abs(float64(newCount) - float64(count)))
	scaleNeeded := (diff != 0)

	if !scaleNeeded {
		s.logger.Info().Str("job", jobName).Msg("No scaling needed/ possible.")
		return sokar.ScaleResult{
			State:            sokar.ScaleIgnored,
			StateDescription: "No scaling needed/ possible.",
			NewCount:         newCount,
		}
	}

	scaleTypeStr := amountToScaleType(diff)
	s.logger.Info().Str("job", jobName).Msgf("Scale %s by %d to %d.", scaleTypeStr, diff, newCount)

	// Set the new job count
	err = s.scalingTarget.SetJobCount(s.job.jobName, newCount)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error adjusting job count to %d: %s.", newCount, err.Error()),
		}
	}

	return sokar.ScaleResult{
		State:            sokar.ScaleDone,
		StateDescription: "Scaling successfully done.",
		NewCount:         newCount,
	}
}

// ScaleBy Scales the target component by the given amount of instances
func (s *Scaler) ScaleBy(amount int) sokar.ScaleResult {
	if r, ok := trueIfNil(s); ok {
		return r
	}

	jobName := s.job.jobName
	min := s.job.minCount
	max := s.job.maxCount
	scaleTypeStr := amountToScaleType(amount)

	s.logger.Info().Str("job", jobName).Msgf("Request to scale job %s by %d.", scaleTypeStr, amount)

	dead, err := s.scalingTarget.IsJobDead(jobName)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error obtaining if job is dead: %s.", err.Error()),
		}
	}

	if dead {
		return sokar.ScaleResult{
			State:            sokar.ScaleIgnored,
			StateDescription: fmt.Sprintf("Job '%s' is dead. Can't scale", jobName),
		}
	}

	count, err := s.scalingTarget.GetJobCount(jobName)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error obtaining job count: %s.", err.Error()),
		}
	}

	chkResult := checkScalingPolicy(count, amount, min, max)
	newCount := chkResult.validCount
	if chkResult.minPolicyViolated {
		s.logger.Info().Str("job", jobName).Msgf("Job.MinCount (%d) policy violated (wanted %d, have %d). Scale %s limited to %d.", min, chkResult.desiredCount, count, scaleTypeStr, newCount)
	}
	if chkResult.maxPolicyViolated {
		s.logger.Info().Str("job", jobName).Msgf("Job.MaxCount (%d) policy violated (wanted %d, have %d). Scale %s limited to %d.", max, chkResult.desiredCount, count, scaleTypeStr, newCount)
	}

	diff := int(math.Abs(float64(newCount) - float64(count)))
	scaleNeeded := (diff != 0)

	if !scaleNeeded {
		s.logger.Info().Str("job", jobName).Msg("No scaling needed/ possible.")
		return sokar.ScaleResult{
			State:            sokar.ScaleIgnored,
			StateDescription: "No scaling needed/ possible.",
			NewCount:         newCount,
		}
	}

	s.logger.Info().Str("job", jobName).Msgf("Scale %s by %d to %d.", scaleTypeStr, diff, newCount)

	// Set the new job count
	err = s.scalingTarget.SetJobCount(s.job.jobName, newCount)
	if err != nil {
		return sokar.ScaleResult{
			State:            sokar.ScaleFailed,
			StateDescription: fmt.Sprintf("Error adjusting job count from %d to %d: %s.", count, newCount, err.Error()),
		}
	}

	return sokar.ScaleResult{
		State:            sokar.ScaleDone,
		StateDescription: "Scaling successfully done.",
		NewCount:         newCount,
	}
}
