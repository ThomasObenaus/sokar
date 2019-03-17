package sokar

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	m "github.com/thomasobenaus/sokar/metrics"
)

// Metrics represents the collection of metrics internally set by sokar.
type Metrics struct {
	scaleEventsTotal   m.Counter
	failedScalingTotal m.Counter
	plannedCount       m.Gauge
}

// NewMetrics returns the metrics collection needed for the SAA.
func NewMetrics() Metrics {

	scaleEventsTotal := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "sokar",
		Name:      "scale_events_total",
		Help:      "Number of received ScaleEvents in total.",
	})

	failedScalingTotal := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "sokar",
		Name:      "failed_scaling_total",
		Help:      "Number of failed scaling actions in total.",
	})

	plannedCount := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Name:      "planned_count",
		Help:      "The count planned by the CapacityPlanner for the current scale action.",
	})

	return Metrics{
		scaleEventsTotal:   scaleEventsTotal,
		failedScalingTotal: failedScalingTotal,
		plannedCount:       plannedCount,
	}
}
