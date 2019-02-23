package alertmanager

import (
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
	c.logger.Info().Msg("SSSSSSSSSSSSSSSSSSS")
}
