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
	currentCount       m.Gauge
	scaleFactor        m.Gauge
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
		Subsystem: "cap",
		Name:      "planned_count",
		Help:      "The count planned by the CapacityPlanner for the current scale action.",
	})

	currentCount := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "cap",
		Name:      "current_count",
		Help:      "The count currently deployed. Based on this count sokar does the planning.",
	})

	scaleFactor := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "cap",
		Name:      "scale_factor",
		Help:      "The scale factor (gradient) as it was received with a ScalingEvent.",
	})

	return Metrics{
		scaleEventsTotal:   scaleEventsTotal,
		failedScalingTotal: failedScalingTotal,
		plannedCount:       plannedCount,
		currentCount:       currentCount,
		scaleFactor:        scaleFactor,
	}
}
