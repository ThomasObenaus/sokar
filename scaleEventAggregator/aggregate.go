package scaleEventAggregator

func (sc *ScaleEventAggregator) aggregate() {
	sc.logger.Info().Msg("Aggregate")
	sc.logPool()

	scaleUp := 0
	sc.scaleAlertPool.iterate(func(key uint32, entry ScaleAlertPoolEntry) {
		scaleUp++
	})

	// FIXME: This currently blocks until the deployment is done
	if scaleUp > 0 {
		sc.logger.Info().Msg("Scale UP by 1")
		sc.emitScaleEvent(1)
	} else {
		sc.logger.Info().Msg("Scale DOWN by 1")
		sc.emitScaleEvent(-1)
	}
}

func (sc *ScaleEventAggregator) logPool() {
	sc.logger.Debug().Int("num-entries", sc.scaleAlertPool.size()).Msg("ScaleAlertPool:")

	sc.scaleAlertPool.iterate(func(key uint32, entry ScaleAlertPoolEntry) {
		sc.logger.Debug().Str("name", entry.scaleAlert.Name).Str("receiver", entry.receiver).Msgf("[%d] fire=%t,start=%s,exp=%s", key, entry.scaleAlert.Firing, entry.scaleAlert.StartedAt.String(), entry.expiresAt.String())
	})
}
