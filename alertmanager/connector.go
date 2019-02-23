package alertmanager

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/scaleEventAggregator"
)

// Connector is the integration of prometheus/alertmanager
type Connector struct {
	logger        zerolog.Logger
	subscriptions []chan scaleEventAggregator.ScaleAlertList
}

// Config cfg for the connector
type Config struct {
	Logger zerolog.Logger
}

// New creates a new instance of the prometheus/alertmanager Connector
func (cfg Config) New() *Connector {
	return &Connector{
		logger: cfg.Logger,
	}
}

// Subscribe is used to register/ subscribe for the channel where scaling alerts are emitted
func (c *Connector) Subscribe(subscriber chan scaleEventAggregator.ScaleAlertList) {
	c.subscriptions = append(c.subscriptions, subscriber)
}

func (c *Connector) fireScaleAlert(scalingAlerts scaleEventAggregator.ScaleAlertList) {
	for _, subscriber := range c.subscriptions {
		subscriber <- scalingAlerts
	}
}

// HandleScaleAlerts is the http end-point implementation for receiving alerts from alertmanager
func (c *Connector) HandleScaleAlerts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.logger.Info().Msg("Receiving scaling alerts")

	defer r.Body.Close()

	alertmanagerResponse := response{}
	err := json.NewDecoder(r.Body).Decode(&alertmanagerResponse)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse data received from alertmanager: %s.", err)
		c.logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	c.logger.Info().Msgf("%d Scaling Alerts received. Will send them to the subscriber.", len(alertmanagerResponse.Alerts))
	scalingAlerts := amResponseToScalingAlerts(alertmanagerResponse)
	c.fireScaleAlert(scalingAlerts)

	w.WriteHeader(http.StatusOK)
}
