package scaleEventAggregator

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/thomasobenaus/sokar/sokar"
)

// Subscribe is used to register for receiving ScaleEvents
func (sc *ScaleEventAggregator) Subscribe(subscriber chan sokar.ScaleEvent) {
	sc.subscriptions = append(sc.subscriptions, subscriber)
}

// ScaleEvent implements the http end-point for emitting a scaling event
func (sc *ScaleEventAggregator) ScaleEvent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	sc.logger.Info().Msg("ScaleEvent Received")
	sc.emitScaleEvent(1)
	w.WriteHeader(http.StatusOK)
}

func (sc *ScaleEventAggregator) emitScaleEvent(scaleFactor float32) {

	for _, subscriber := range sc.subscriptions {
		subscriber <- sokar.ScaleEvent{ScaleFactor: scaleFactor}
	}
}

// Run starts the ScaleEventAggregator
func (sc *ScaleEventAggregator) Run() {

	sc.logger.Info().Msg("Subscribe at scale alert receivers")
	scaleAlertChannel := make(chan ScaleAlertList)
	for _, receiver := range sc.receivers {
		receiver.Subscribe(scaleAlertChannel)
	}

	aggregationTicker := time.NewTicker(5000 * time.Millisecond)
	cleanupTicker := time.NewTicker(20000 * time.Millisecond)

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

func (sc *ScaleEventAggregator) handleReceivedScaleAlerts(scaleAlerts ScaleAlertList) {
	sc.logger.Info().Msgf("%d Alerts received.", len(scaleAlerts))
	sc.scaleAlertPool.update(scaleAlerts)
}

// Stop tears down ScaleEventAggregator
func (sc *ScaleEventAggregator) Stop() {
	sc.logger.Info().Msg("Teardown requested")
	// send the stop message
	sc.stopChan <- struct{}{}
}

// Join blocks/ waits until ScaleEventAggregator has been stopped
func (sc *ScaleEventAggregator) Join() {
	<-sc.stopChan
}
