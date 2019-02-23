package scaleEventAggregator

import (
	"net/http"

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

	for _, subscribor := range sc.subscriptions {
		subscribor <- sokar.ScaleEvent{ScaleFactor: 1}
	}
}
