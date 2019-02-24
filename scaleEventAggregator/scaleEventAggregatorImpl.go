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
	sc.emitScaleEvent()
	w.WriteHeader(http.StatusOK)
}

func (sc *ScaleEventAggregator) emitScaleEvent() {

	for _, subscriber := range sc.subscriptions {
		subscriber <- sokar.ScaleEvent{ScaleFactor: 1}
	}
}

func (sc *ScaleEventAggregator) evaluate() {
	sc.logger.Info().Msg("Evaluate")

	//for alertName := range sc.alertMap {
	//
	//	sf, ok := sc.scaleFactorMap[alertName]
	//	if !ok {
	//		log.Fatal("Alert not in map")
	//	}
	//
	//	sc.logger.Info().Msgf("Alert %s, SF %f", alertName, sf)
	//
	//	// HACK: BLOCKS UNTIL DEPLOYMENT IS DONE!!!!
	//	sc.emitScaleEvent()
	//}

}

// Run starts the ScaleEventAggregator
func (sc *ScaleEventAggregator) Run() {

	sc.logger.Info().Msg("Subscribe at scale alert receivers")
	scaleAlertChannel := make(chan ScaleAlertList)
	for _, receiver := range sc.receivers {
		receiver.Subscribe(scaleAlertChannel)
	}

	ticker := time.NewTicker(1000 * time.Millisecond)

	// main loop
	go func() {
		sc.logger.Info().Msg("Main process loop started")

	loop:
		for {
			select {
			case <-sc.stopChan:
				ticker.Stop()
				close(sc.stopChan)
				break loop

			case <-ticker.C:
				sc.evaluate()
			case scaleAlerts := <-scaleAlertChannel:
				sc.logger.Info().Msgf("SCCCCCCCCCCCCCALE %+v", scaleAlerts)

				//for _, alert := range scaleAlerts {
				//	if !alert.Firing {
				//		delete(sc.alertMap, alert.Name)
				//	} else {
				//		sc.alertMap[alert.Name] = alert
				//	}
				//
				//}

			}
		}
		sc.logger.Info().Msg("Main process loop left")

	}()

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
