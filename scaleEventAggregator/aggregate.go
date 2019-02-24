package scaleEventAggregator

func (sc *ScaleEventAggregator) aggregate() {
	sc.logger.Info().Msg("Aggregate")
	sc.logPool()

	scaleUp := 0
	sc.scaleAlertPool.iterate(func(alert ScaleAlert) {
		scaleUp++
	})

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

	sc.scaleAlertPool.iterate(func(scaleAlert ScaleAlert) {
		sc.logger.Debug().Str("name", scaleAlert.Name).Bool("fires", scaleAlert.Firing).Str("startedAt", scaleAlert.StartedAt.String()).Msg("\t")
	})
}
