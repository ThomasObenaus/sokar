package scaler

import (
	"fmt"
	"time"

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

// scale scales the scalingObject from currentCount to desiredCount.
// Internally it is checked if a scaling is needed and if the scaling policy is valid.
// If the force flag is true then even in dry-run mode the scaling will be applied.
func (s *Scaler) scale(desiredCount uint, currentCount uint, force bool) scaleResult {

	sObjName := s.scalingObject.Name
	min := s.scalingObject.MinCount
	max := s.scalingObject.MaxCount

	s.logger.Info().Str("scalingObject", sObjName).Msgf("Request to scale scalingObject from %d to %d.", currentCount, desiredCount)

	dead, err := s.scalingTarget.IsScalingObjectDead(sObjName)
	if err != nil {
		return scaleResult{
			state:            scaleFailed,
			stateDescription: fmt.Sprintf("Error obtaining if scalingObject is dead: %s.", err.Error()),
			newCount:         currentCount,
		}
	}

	if dead {
		return scaleResult{
			state:            scaleIgnored,
			stateDescription: fmt.Sprintf("ScalingObject '%s' is dead. Can't scale", sObjName),
			newCount:         currentCount,
		}
	}

	chkResult := checkScalingPolicy(desiredCount, min, max)
	newCount := chkResult.validCount
	if chkResult.minPolicyViolated {
		s.logger.Info().Str("scalingObject", sObjName).Msgf("ScalingObject.MinCount (%d) policy violated (wanted %d, have %d). Scale limited to %d.", min, chkResult.desiredCount, currentCount, newCount)
		s.metrics.scalingPolicyViolated.WithLabelValues("min").Inc()
	}
	if chkResult.maxPolicyViolated {
		s.logger.Info().Str("scalingObject", sObjName).Msgf("ScalingObject.MaxCount (%d) policy violated (wanted %d, have %d). Scale limited to %d.", max, chkResult.desiredCount, currentCount, newCount)
		s.metrics.scalingPolicyViolated.WithLabelValues("max").Inc()
	}

	diff := helper.SubUint(newCount, currentCount)
	scaleNeeded := (diff != 0)

	if !scaleNeeded {
		s.logger.Info().Str("scalingObject", sObjName).Msg("No scaling needed/ possible.")
		return scaleResult{
			state:            scaleIgnored,
			stateDescription: "No scaling needed/ possible.",
			newCount:         newCount,
		}
	}

	// memorize the time the scaling started
	if scaleNeeded && isScalePermitted(s.dryRunMode, force) {
		s.lastScaleAction = time.Now()
	}

	// TODO: Go in cooldown just incase the scaling was successful
	// TODO: Split cooldown for up and downscaling
	// See: https://github.com/ThomasObenaus/sokar/issues/111
	return s.executeScale(currentCount, newCount, force)
}

// executeScale just executes the wanted scale, even if it is not needed in case currentCount == newCount.
// The only thing that is checked is if the scaler is in dry run mode or not.
// In dry-run mode the scaling is not applied by actually scaling the scalingObject, only a metric is
// updated, that reflects the fact that a scaling was skipped/ ignored.
// Only the force flag can be used to override this behavior. If the force flag is true then even in
// dry-run mode the scaling will be applied.
func (s *Scaler) executeScale(currentCount, newCount uint, force bool) scaleResult {
	sObjName := s.scalingObject.Name
	min := s.scalingObject.MinCount
	max := s.scalingObject.MaxCount

	diff := helper.SubUint(newCount, currentCount)
	scaleTypeStr := amountToScaleType(diff)

	// the force flag can overrule the dry run mode
	if !isScalePermitted(s.dryRunMode, force) {
		s.logger.Info().Str("scalingObject", sObjName).Msgf("Skip scale %s by %d to %d (DryRun, force=%v).", scaleTypeStr, diff, newCount, force)
		s.metrics.plannedButSkippedScalingOpen.WithLabelValues(scaleTypeStr).Set(1)

		return scaleResult{
			state:            scaleIgnored,
			stateDescription: "Scaling skipped - dry run mode.",
			newCount:         currentCount,
		}
	}

	s.logger.Info().Str("scalingObject", sObjName).Msgf("Scale %s by %d to %d (force=%v).", scaleTypeStr, diff, newCount, force)
	s.metrics.plannedButSkippedScalingOpen.WithLabelValues(scaleTypeStr).Set(0)

	err := s.scalingTarget.AdjustScalingObjectCount(sObjName, min, max, currentCount, newCount)
	if err != nil {
		return scaleResult{
			state:            scaleFailed,
			stateDescription: fmt.Sprintf("Error adjusting scalingObject count to %d: %s.", newCount, err.Error()),
			newCount:         currentCount,
		}
	}

	return scaleResult{
		state:            scaleDone,
		stateDescription: "Scaling successfully done.",
		newCount:         newCount,
	}
}

func isScalePermitted(dryRun, force bool) bool {
	return !dryRun || force
}
