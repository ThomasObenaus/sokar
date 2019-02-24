package scaleEventAggregator

func (sc *ScaleEventAggregator) aggregate() {
	sc.logger.Info().Msg("Aggregate")
	sc.logPool()

	sc.scaleAlertPool.iterate(sc.updateScaleCounter)

	// FIXME: This currently blocks until the deployment is done
	if sc.scaleCounter > sc.upScalingThreshold {
		sc.logger.Info().Msgf("Scale UP by 1 because upscalingThreshold (%f) was violated. ScaleCounter is currently %f", sc.upScalingThreshold, sc.scaleCounter)
		sc.scaleCounter = 0
		sc.emitScaleEvent(1)
	} else if sc.scaleCounter < sc.downScalingThreshold {
		sc.logger.Info().Msgf("Scale DOWN by 1 because downScalingThreshold (%f) was violated. ScaleCounter is currently %f", sc.downScalingThreshold, sc.scaleCounter)
		sc.scaleCounter = 0
		sc.emitScaleEvent(-1)
	} else {
		sc.logger.Info().Msgf("No scaling needed. ScaleCounter is currently %f [%f/%f].", sc.scaleCounter, sc.downScalingThreshold, sc.upScalingThreshold)
	}
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
	weight := getWeight(alertName, sc.weightMap)
	sc.scaleCounter += weight

	sc.logger.Debug().Msgf("ScaleCounter updated by %f to %f because of scaling-alert %s", weight, sc.scaleCounter, alertName)
}
