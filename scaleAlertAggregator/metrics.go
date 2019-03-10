package scaleAlertAggregator

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	m "github.com/thomasobenaus/sokar/metrics"
)

// Metrics represents the collection of metrics internally set by
// the ScaleAlertAggregator.
type Metrics struct {
	scaleCounter m.Gauge

	alerts m.GaugeVec
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

	return Metrics{
		scaleCounter: scaleCounter,
		alerts:       alerts,
	}
}
