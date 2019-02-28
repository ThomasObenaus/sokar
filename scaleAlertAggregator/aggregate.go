package scaleAlertAggregator

import (
	"time"
)

// aggregate all ScaleAlerts available and updates internally the scaleCounter.
func (sc *ScaleAlertAggregator) aggregate() {
	sc.logger.Info().Msg("Aggregation")
	sc.logPool()

	// collect all alerts from ScaleAlertPool
	var poolEntries []ScaleAlertPoolEntry
	sc.scaleAlertPool.iterate(func(key uint32, entry ScaleAlertPoolEntry) {
		poolEntries = append(poolEntries, entry)
	})

	alertsChangedScaleCounter := sc.applyAlertsToScaleCounter(poolEntries, sc.weightMap, sc.evaluationCycle)

	if !alertsChangedScaleCounter {
		sc.applyScaleCounterDamping(sc.noAlertScaleDamping, sc.evaluationCycle)
	}
}

// applyScaleCounterDamping applies the given damping to the scaleCounter
func (sc *ScaleAlertAggregator) applyScaleCounterDamping(noAlertScaleDamping float32, evaluationCycle time.Duration) {
	weight := weightPerSecondToWeight(noAlertScaleDamping, evaluationCycle)
	scaleIncrement := computeScaleCounterDamping(sc.scaleCounter, weight)
	sc.scaleCounter += scaleIncrement

	if scaleIncrement != 0 {
		sc.logger.Debug().Msgf("ScaleCounter updated/damped by %f to %f because no scaling-alert changed the scale counter. Damping (per s): %f.", scaleIncrement, sc.scaleCounter, sc.noAlertScaleDamping)
	}
}

// computeScaleCounterDamping computes the value that has to be added to the scaleCounter
// in order to move it more to 0. It is either a positive or negative version of the given dampingFactor.
func computeScaleCounterDamping(scaleCounter float32, dampingFactor float32) float32 {
	negativeDamping := true
	abs := scaleCounter
	if abs < 0 {
		abs = scaleCounter * -1
		negativeDamping = false
	}

	var result float32
	if abs <= dampingFactor {
		result = abs
	} else {
		result = dampingFactor
	}

	if negativeDamping {
		result *= -1
	}

	return result
}

func (sc *ScaleAlertAggregator) logPool() {
	sc.logger.Debug().Int("num-entries", sc.scaleAlertPool.size()).Msg("ScaleAlertPool:")

	sc.scaleAlertPool.iterate(func(key uint32, entry ScaleAlertPoolEntry) {
		sc.logger.Debug().Str("name", entry.scaleAlert.Name).Str("receiver", entry.receiver).Msgf("[%d] fire=%t,start=%s,exp=%s", key, entry.scaleAlert.Firing, entry.scaleAlert.StartedAt.String(), entry.expiresAt.String())
	})
}

// computeScaleCounterIncrement determines how much the scaleCounter has to be changed for the given alert.
func computeScaleCounterIncrement(alertName string, weightMap ScaleAlertWeightMap, evaluationCycle time.Duration) (scaleIncrement float32, weightPerSecond float32) {
	weightPerSecond = getWeight(alertName, weightMap)
	if weightPerSecond == 0 {
		return 0, 0
	}
	scaleIncrement = weightPerSecondToWeight(weightPerSecond, evaluationCycle)
	return scaleIncrement, weightPerSecond
}

// applyAlertsToScaleCounter applies the given alerts to the scaleCounter by incrementing/ decrementing the counter accordingly.
func (sc *ScaleAlertAggregator) applyAlertsToScaleCounter(entries []ScaleAlertPoolEntry, weightMap ScaleAlertWeightMap, evaluationCycle time.Duration) (scaleCounterHasChanged bool) {
	oldScaleCounterValue := sc.scaleCounter

	for _, entry := range entries {
		// ignore resolved alerts
		if !entry.scaleAlert.Firing {
			continue
		}

		alertName := entry.scaleAlert.Name
		scaleIncrement, weightPerSecond := computeScaleCounterIncrement(alertName, weightMap, evaluationCycle)
		sc.scaleCounter += scaleIncrement

		sc.logger.Debug().Msgf("ScaleCounter updated by %f to %f. Scaling-Alert: '%s' (%f wps).", scaleIncrement, sc.scaleCounter, alertName, weightPerSecond)
	}

	return oldScaleCounterValue != sc.scaleCounter
}
