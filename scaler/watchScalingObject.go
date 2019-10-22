package scaler

import "strconv"

// countMeetsExpectations returns false in case the current count does not match the
// expectations. This could either be the case if the current deployed amount
// of allocations is "out of bounds" regarding the defined min/ max of the scalingObject.
// Or if a the current count does not correlate to the desired amount of allocations
// being deployed.
func countMeetsExpectations(current uint, min uint, max uint, desired optionalValue) (asExpected bool, expectedCount uint) {

	asExpected = true
	expectedCount = current

	if desired.isKnown {
		expectedCount = desired.value
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

func (s *Scaler) ensureScalingObjectCount() error {

	count, err := s.GetCount()
	if err != nil {
		return err
	}

	asExpected, expected := countMeetsExpectations(count, s.scalingObject.MinCount, s.scalingObject.MaxCount, s.desiredScale)
	if !asExpected {
		s.logger.Warn().Bool("watcher", true).Msgf("The scalingObject count (%d) is not as expected. Thus the scalingObject has to be rescaled to %d.", count, expected)
		if err := s.openScalingTicket(expected, false); err != nil {
			return err
		}
	} else {
		desiredStr := "n/a"
		if s.desiredScale.isKnown {
			desiredStr = strconv.Itoa(int(s.desiredScale.value))
		}
		s.logger.Debug().Bool("watcher", true).Uint("count", count).Str("desired", desiredStr).Uint("expected", expected).Msg("Count as expected, no adjustment needed.")
	}

	return nil
}
