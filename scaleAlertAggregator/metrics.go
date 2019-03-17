package scaleAlertAggregator

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	m "github.com/thomasobenaus/sokar/metrics"
)

// Metrics represents the collection of metrics internally set by
// the ScaleAlertAggregator.
type Metrics struct {
	scaleCounter      m.Gauge
	alerts            m.GaugeVec
	scaleFactor       m.Gauge
	scaleEventCounter m.CounterVec
}

// NewMetrics returns the metrics collection needed for the SAA.
func NewMetrics() Metrics {

	scaleCounter := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "saa",
		Name:      "scale_counter",
		Help:      "The current value of the ScaleCounter. This is the aggregated weights of all ScalingAlerts.",
	})

	alertLabels := []string{"direction"}
	alerts := m.NewWrappedGaugeVec(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "saa",
		Name:      "alerts_total",
		Help:      "The number of currently active down and up alerts.",
	}, alertLabels)

	scaleFactor := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "saa",
		Name:      "scale_factor",
		Help:      "The scale factor (gradient) as it is calculated by the SAA on each evaluation.",
	})

	scaleEventDirs := []string{"direction"}
	scaleEventCounter := m.NewWrappedCounterVec(prometheus.CounterOpts{
		Namespace: "sokar",
		Subsystem: "saa",
		Name:      "scale_event_counter",
		Help:      "Counts the number of ScaleEvent's separated by up/down.",
	}, scaleEventDirs)

	return Metrics{
		scaleCounter:      scaleCounter,
		alerts:            alerts,
		scaleFactor:       scaleFactor,
		scaleEventCounter: scaleEventCounter,
	}
}
