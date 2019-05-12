package scaler

import "strconv"

// countMeetsExpectations returns false in case the current count does not match the
// expectations. This could either be the case if the current deployed amount
// of allocations is "out of bounds" regarding the defined min/ max of the job.
// Or if a the current count does not correlate to the desired amount of allocations
// being deployed.
func countMeetsExpectations(current uint, min uint, max uint, desired *uint) (asExpected bool, expectedCount uint) {

	asExpected = true
	expectedCount = current

	if desired != nil {
		expectedCount = *desired
	}

	if current != expectedCount {
		asExpected = false
	}

	if current < min {
		expectedCount = min
		asExpected = false
	}

	if current > max {
		expectedCount = max
		asExpected = false
	}

	return asExpected, expectedCount
}

func (s *Scaler) ensureJobCount() error {

	count, err := s.GetCount()
	if err != nil {
		return err
	}

	asExpected, expected := countMeetsExpectations(count, s.job.minCount, s.job.maxCount, s.desiredScale)
	if !asExpected {
		s.logger.Warn().Msgf("The job count (%d) was not as expected. Thus the job had to be rescaled to %d.", count, expected)
		if err := s.openScalingTicket(expected, false); err != nil {
			return err
		}
	} else {
		desiredStr := "n/a"
		if s.desiredScale != nil {
			desiredStr = strconv.Itoa(int(*s.desiredScale))
		}
		s.logger.Debug().Uint("count", count).Str("desired", desiredStr).Uint("expected", expected).Msg("Count as expected, no adjustment needed.")
	}

	return nil
}
