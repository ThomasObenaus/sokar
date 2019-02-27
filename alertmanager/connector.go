package alertmanager

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/rs/zerolog"
	saa "github.com/thomasobenaus/sokar/scaleAlertAggregator"
)

// Connector is the integration of prometheus/alertmanager
type Connector struct {
	logger zerolog.Logger

	// handleFuncs is a list of registered handlers for received ScaleAlerts
	handleFuncs []saa.ScaleAlertHandleFunc
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

// Register is used to register the given handler func.
// The ScaleAlertHandleFunc is called each time the alertmanager connector receives an alert.
func (c *Connector) Register(handleFunc saa.ScaleAlertHandleFunc) {
	c.handleFuncs = append(c.handleFuncs, handleFunc)
}

// fireScaleAlertPacket sends the given ScaleAlertPacket to all registered handler functions.
func (c *Connector) fireScaleAlertPacket(scalingAlerts saa.ScaleAlertPacket) {
	for _, handleFunc := range c.handleFuncs {
		handleFunc(scalingAlerts.Emitter, scalingAlerts)
	}
}

// HandleScaleAlerts is the http end-point implementation for receiving alerts from alertmanager
func (c *Connector) HandleScaleAlerts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.logger.Info().Msg("Received scaling alert packet.")

	defer r.Body.Close()

	alertmanagerResponse := response{}
	err := json.NewDecoder(r.Body).Decode(&alertmanagerResponse)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse data received from alertmanager: %s.", err)
		c.logger.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	scalingAlertPacket := amResponseToScalingAlerts(alertmanagerResponse)
	c.logger.Info().Msgf("%d Scaling Alerts received from '%s'. Will send them to the subscriber.", len(scalingAlertPacket.ScaleAlerts), scalingAlertPacket.Emitter)
	c.fireScaleAlertPacket(scalingAlertPacket)

	w.WriteHeader(http.StatusOK)
}
