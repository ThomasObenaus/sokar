package scaleAlertAggregator

import (
	"time"

	"github.com/thomasobenaus/sokar/sokar/iface"
)

// Subscribe is used to register for receiving ScaleEvents
func (sc *ScaleAlertAggregator) Subscribe(subscriber chan sokar.ScaleEvent) {
	sc.subscriptions = append(sc.subscriptions, subscriber)
}

func (sc *ScaleAlertAggregator) emitScaleEvent(scaleFactor float32) {

	for _, subscriber := range sc.subscriptions {
		subscriber <- sokar.ScaleEvent{ScaleFactor: scaleFactor}
	}
}

// GetName returns the name of this component
func (sc *ScaleAlertAggregator) GetName() string {
	return "scaleAlertAggregator"
}

// Run starts the ScaleAlertAggregator
func (sc *ScaleAlertAggregator) Run() {

	sc.logger.Info().Msg("Register at scale alert emitters")
	for _, emitter := range sc.emitters {
		emitter.Register(sc.handleScaleAlerts)
	}

	evaluationTicker := time.NewTicker(sc.evaluationCycle)
	cleanupTicker := time.NewTicker(sc.cleanupCycle)

	// main loop
	go func() {
		sc.logger.Info().Msg("Main process loop started")

	loop:
		for {
			select {

			case <-sc.stopChan:
				evaluationTicker.Stop()
				cleanupTicker.Stop()
				close(sc.stopChan)
				break loop

			case <-cleanupTicker.C:
				sc.scaleAlertPool.cleanup()

			case <-evaluationTicker.C:
				sc.aggregate()
				gradient := sc.evaluate()

				if gradient != 0 {
					sc.emitScaleEvent(gradient)
				}
			}
		}
		sc.logger.Info().Msg("Main process loop left")
	}()

}

func (sc *ScaleAlertAggregator) handleScaleAlerts(emitter string, scaPckg ScaleAlertPacket) {
	sc.logger.Info().Msgf("%d Alerts received from %s.", len(scaPckg.ScaleAlerts), emitter)
	sc.scaleAlertPool.update(emitter, scaPckg.ScaleAlerts, sc.weightMap)

	updateAlertMetrics(&sc.scaleAlertPool, &sc.metrics)
	sc.logPool()
}

func updateAlertMetrics(pool *ScaleAlertPool, metrics *Metrics) {
	numUp := float64(0)
	numDown := float64(0)
	numNeutral := float64(0)

	pool.iterate(func(key uint32, entry ScaleAlertPoolEntry) {
		if entry.weight > 0 {
			numUp++
		} else if entry.weight < 0 {
			numDown++
		} else {
			numNeutral++
		}
	})

	metrics.alerts.WithLabelValues("up").Set(numUp)
	metrics.alerts.WithLabelValues("down").Set(numDown)
	metrics.alerts.WithLabelValues("neutral").Set(numNeutral)
}

// Stop tears down ScaleAlertAggregator
func (sc *ScaleAlertAggregator) Stop() {
	sc.logger.Info().Msg("Teardown requested")
	// send the stop message
	sc.stopChan <- struct{}{}
}

// Join blocks/ waits until ScaleAlertAggregator has been stopped
func (sc *ScaleAlertAggregator) Join() {
	<-sc.stopChan
}
