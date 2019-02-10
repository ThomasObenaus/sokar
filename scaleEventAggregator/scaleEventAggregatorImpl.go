package scaleEventAggregator

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thomasobenaus/sokar/sokar"
)

func (sc *ScaleEventAggregator) Substribe(subscriber chan sokar.ScaleEvent) {
	sc.subscriptions = append(sc.subscriptions, subscriber)
}

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
