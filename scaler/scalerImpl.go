package scaler

import (
	"fmt"
	"math"
)

type jobConfig struct {
	jobName  string
	minCount uint
	maxCount uint
}

// ScaleBy Scales the target component by the given amount of instances
func (s *Scaler) ScaleBy(amount int) error {
	jobName := s.job.jobName

	scaleTypeStr := "UP"
	if amount < 0 {
		scaleTypeStr = "DOWN"
	}
	s.logger.Info().Str("job", jobName).Msgf("Request to scale job %s by %d.", scaleTypeStr, amount)

	dead, err := s.scalingTarget.IsJobDead(jobName)
	if err != nil {
		return fmt.Errorf("Error obtaining if job is dead: %s.", err.Error())
	}

	if dead {
		return fmt.Errorf("Job '%s' is dead. Can't scale", jobName)
	}

	count, err := s.scalingTarget.GetJobCount(jobName)
	if err != nil {
		return fmt.Errorf("Error obtaining job count: %s.", err.Error())
	}

	newCountTmp := int(count) + amount
	newCount := uint(newCountTmp)
	if newCountTmp < 0 {
		newCount = 0
	}
	countWanted := newCount

	// check if count exceeds minimum
	if newCount < s.job.minCount {
		newCount = s.job.minCount
		s.logger.Info().Str("job", jobName).Msgf("Job.MinCount (%d) policy violated (wanted %d, have %d). Scale %s limited to %d.", s.job.minCount, countWanted, count, scaleTypeStr, newCount)
	}

	// check if count exceeds maximum
	if newCount > s.job.maxCount {
		newCount = s.job.maxCount
		s.logger.Info().Str("job", jobName).Msgf("Job.MaxCount (%d) policy violated (wanted %d, have %d). Scale %s limited to %d.", s.job.maxCount, countWanted, count, scaleTypeStr, newCount)
	}

	diff := int(math.Abs(float64(newCount) - float64(count)))
	scaleNeeded := (diff != 0)

	if !scaleNeeded {
		s.logger.Info().Str("job", jobName).Msg("No scaling needed/ possible.")
		return nil
	}

	s.logger.Info().Str("job", jobName).Msgf("Scale %s by %d to %d.", scaleTypeStr, diff, newCount)

	// Set the new job count
	err = s.scalingTarget.SetJobCount(s.job.jobName, newCount)
	if err != nil {
		return fmt.Errorf("Error adjusting job count from %d to %d: %s.", count, newCount, err.Error())
	}

	return nil
}
