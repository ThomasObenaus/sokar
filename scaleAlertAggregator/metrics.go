package scaleAlertAggregator

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	m "github.com/thomasobenaus/sokar/metrics"
)

type Metrics struct {
	scaleCounter m.Gauge
}

// NewMetrics returns the metrics collection needed for the SAA.
func NewMetrics() Metrics {

	scaleCounter := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "saa",
		Name:      "scale_counter",
		Help:      "The current value of the ScaleCounter. This is the aggregated weights of all ScalingAlerts.",
	})

	return Metrics{
		scaleCounter: scaleCounter,
	}
}
