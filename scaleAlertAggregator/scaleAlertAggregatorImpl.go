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
