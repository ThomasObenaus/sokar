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

	sc.logger.Info().Msg("Subscribe at scale alert receivers")
	scaleAlertChannel := make(chan ScaleAlertPacket)
	for _, receiver := range sc.receivers {
		receiver.Subscribe(scaleAlertChannel)
	}

	aggregationTicker := time.NewTicker(sc.aggregationCycle)
	cleanupTicker := time.NewTicker(sc.cleanupCycle)

	// main loop
	go func() {
		sc.logger.Info().Msg("Main process loop started")

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
				sc.aggregate()

			case scaleAlerts := <-scaleAlertChannel:
				sc.handleReceivedScaleAlerts(scaleAlerts)
			}
		}
		sc.logger.Info().Msg("Main process loop left")
	}()

}

func (sc *ScaleAlertAggregator) handleReceivedScaleAlerts(scaPckg ScaleAlertPacket) {
	sc.logger.Info().Msgf("%d Alerts received from %s.", len(scaPckg.ScaleAlerts), scaPckg.Emitter)
	sc.scaleAlertPool.update(scaPckg.Emitter, scaPckg.ScaleAlerts)
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
