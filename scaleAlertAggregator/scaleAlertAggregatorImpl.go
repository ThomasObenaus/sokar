package scaleAlertAggregator

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/thomasobenaus/sokar/sokar"
)

// Subscribe is used to register for receiving ScaleEvents
func (sc *ScaleAlertAggregator) Subscribe(subscriber chan sokar.ScaleEvent) {
	sc.subscriptions = append(sc.subscriptions, subscriber)
}

// ScaleEvent implements the http end-point for emitting a scaling event
func (sc *ScaleAlertAggregator) ScaleEvent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	sc.logger.Info().Msg("ScaleEvent Received")
	sc.emitScaleEvent(1)
	w.WriteHeader(http.StatusOK)
}

func (sc *ScaleAlertAggregator) emitScaleEvent(scaleFactor float32) {

	for _, subscriber := range sc.subscriptions {
		subscriber <- sokar.ScaleEvent{ScaleFactor: scaleFactor}
	}
}

// Run starts the ScaleAlertAggregator
func (sc *ScaleAlertAggregator) Run() {

	sc.logger.Info().Msg("Register at scale alert emitters")
	for _, emitter := range sc.emitters {
		emitter.Register(sc.handleScaleAlerts)
	}

	aggregationTicker := time.NewTicker(sc.aggregationCycle)
	cleanupTicker := time.NewTicker(sc.cleanupCycle)

	// main loop
	go func() {
		sc.logger.Info().Msg("Main process loop started")

		aggregationCounter := uint(0)

	loop:
		for {
			select {

			case <-sc.stopChan:
				aggregationTicker.Stop()
				cleanupTicker.Stop()
				close(sc.stopChan)
				break loop

			case <-cleanupTicker.C:
				sc.scaleAlertPool.cleanup()

			case <-aggregationTicker.C:
				aggregationCounter++
				sc.aggregate()

				if sc.isScalingNeeded() {
					gradient := sc.scaleCounterGradient.UpdateAndGet(sc.scaleCounter, time.Now())
					sc.logger.Info().Msgf("Scaling needed. Gradient %f.", gradient)
					// TODO: use this LL somewhere
					//sc.logger.Info().Msgf("Scale by %f because upscalingThreshold (%f) was violated. ScaleCounter is currently %f", scaleFactor, sc.upScalingThreshold, sc.scaleCounter)

					// FIXME: This currently blocks until the deployment is done
					sc.emitScaleEvent(gradient)
					sc.scaleCounter = 0
					sc.scaleCounterGradient.Update(0, time.Now())
					aggregationCounter = 0
				} else if aggregationCounter%sc.evaluationPeriodFactor == 0 {
					gradient := sc.scaleCounterGradient.UpdateAndGet(sc.scaleCounter, time.Now())
					sc.logger.Debug().Msgf("Evaluation period exceeded. Refresh gradient %f.", gradient)
				}

				// TODO: Use this LL somewhere
				//	sc.logger.Info().Msgf("No scaling needed. ScaleCounter is currently %f [%f/%f/%f].", sc.scaleCounter, sc.downScalingThreshold, sc.upScalingThreshold, sc.noAlertScaleDamping)

			}
		}
		sc.logger.Info().Msg("Main process loop left")
	}()

}

func (sc *ScaleAlertAggregator) handleScaleAlerts(emitter string, scaPckg ScaleAlertPacket) {
	sc.logger.Info().Msgf("%d Alerts received from %s.", len(scaPckg.ScaleAlerts), emitter)
	sc.scaleAlertPool.update(emitter, scaPckg.ScaleAlerts)
	sc.logPool()
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
