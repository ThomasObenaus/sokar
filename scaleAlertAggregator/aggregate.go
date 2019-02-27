package scaleEventAggregator

import (
	"time"
)

func (sc *ScaleEventAggregator) aggregate() {
	sc.logger.Info().Msg("Aggregate")
	sc.logPool()

	sc.scaleAlertPool.iterate(sc.updateScaleCounter)

	// FIXME: This currently blocks until the deployment is done
	if sc.scaleCounter > sc.upScalingThreshold {
		sc.logger.Info().Msgf("Scale UP by 1 because upscalingThreshold (%f) was violated. ScaleCounter is currently %f", sc.upScalingThreshold, sc.scaleCounter)
		sc.emitScaleEvent(1)
		sc.scaleCounter = 0
	} else if sc.scaleCounter < sc.downScalingThreshold {
		sc.logger.Info().Msgf("Scale DOWN by 1 because downScalingThreshold (%f) was violated. ScaleCounter is currently %f", sc.downScalingThreshold, sc.scaleCounter)
		sc.emitScaleEvent(-1)
		sc.scaleCounter = 0
	} else {
		sc.logger.Info().Msgf("No scaling needed. ScaleCounter is currently %f [%f/%f/%f].", sc.scaleCounter, sc.downScalingThreshold, sc.upScalingThreshold, sc.noAlertScaleDamping)

		weight := weightPerSecondToWeight(sc.noAlertScaleDamping, sc.aggregationCycle)
		sc.scaleCounter += computeScaleCounterDamping(sc.scaleCounter, weight)
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

func (sc *ScaleEventAggregator) logPool() {
	sc.logger.Debug().Int("num-entries", sc.scaleAlertPool.size()).Msg("ScaleAlertPool:")

	sc.scaleAlertPool.iterate(func(key uint32, entry ScaleAlertPoolEntry) {
		sc.logger.Debug().Str("name", entry.scaleAlert.Name).Str("receiver", entry.receiver).Msgf("[%d] fire=%t,start=%s,exp=%s", key, entry.scaleAlert.Firing, entry.scaleAlert.StartedAt.String(), entry.expiresAt.String())
	})
}

// getWeight returns the scale weight for the given alert.
// 0 is returned in case the weight for the given alert is not defined in the map
func getWeight(alertName string, weightMap ScaleAlertWeightMap) float32 {
	w, ok := weightMap[alertName]
	if !ok {
		return 0
	}
	return w
}

func (sc *ScaleEventAggregator) updateScaleCounter(key uint32, entry ScaleAlertPoolEntry) {
	// ignore resolved alerts
	if !entry.scaleAlert.Firing {
		return
	}

	alertName := entry.scaleAlert.Name
	weightPerSecond := getWeight(alertName, sc.weightMap)
	scaleIncrement := weightPerSecondToWeight(weightPerSecond, sc.aggregationCycle)
	sc.scaleCounter += scaleIncrement

	sc.logger.Debug().Msgf("ScaleCounter updated by %f to %f because of a scaling-alert (name=%s, weight=%f).", scaleIncrement, sc.scaleCounter, alertName, weightPerSecond)
}

//weightPerSecondToWeight converts the given weight (per second) into an absolute weight
// based on the given aggregate cycle.
func weightPerSecondToWeight(weightPerSecond float32, aggregationCycle time.Duration) float32 {
	return float32(aggregationCycle.Seconds() * float64(weightPerSecond))
}
