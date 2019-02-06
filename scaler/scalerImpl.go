package scaler

import (
	"fmt"
	"math"

	"github.com/rs/zerolog"
)

type jobConfig struct {
	jobName  string
	minCount uint
	maxCount uint
}

type scalerImpl struct {
	logger        zerolog.Logger
	scalingTarget ScalingTarget
	job           jobConfig
}

func (s *scalerImpl) ScaleBy(amount int) error {
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
		s.logger.Info().Str("job", jobName).Msg("Job is dead. Makes no sense to scale.")
		return nil
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
