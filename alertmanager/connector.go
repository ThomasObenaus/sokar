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
	subscriptions []chan scaleEventAggregator.ScaleAlert
}

// Config cfg for the connector
type Config struct {
	Logger zerolog.Logger
}

// New creates a new instance of the prometheus/alertmanager Connector
func (cfg Config) New() Connector {
	return Connector{
		logger: cfg.Logger,
	}
}

func (c *Connector) Subscribe(subscriber chan scaleEventAggregator.ScaleAlert) {
	c.subscriptions = append(c.subscriptions, subscriber)
}

func (c *Connector) fireScaleAlert(scaleAlert scaleEventAggregator.ScaleAlert) {
	for _, subscriber := range c.subscriptions {
		subscriber <- scaleAlert
	}
}

func (c *Connector) HandleScaleAlert(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.logger.Info().Msg("Receiving scaling alerts")

	defer r.Body.Close()

	alerts := response{}
	err := json.NewDecoder(r.Body).Decode(&alerts)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse data received from alertmanager: %s.", err)
		c.logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	c.logger.Info().Msgf("%d Scaling Alerts received.", len(alerts.Alerts))
	for _, alert := range alerts.Alerts {
		c.logger.Info().Str("status", alert.Status).Msgf("Labels: %+v", alert.Labels)
	}

	w.WriteHeader(http.StatusOK)
}
