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

// ScaleBy Scales the target component by the given amount of instances
func (s *Scaler) ScaleBy(amount int) sokar.ScaleResult {
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
