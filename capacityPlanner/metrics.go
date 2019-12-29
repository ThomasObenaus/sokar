package capacityPlanner

import (
	"github.com/prometheus/client_golang/prometheus"
	m "github.com/thomasobenaus/sokar/metrics"
)

// Metrics represents the collection of metrics internally set by CapacityPlanner
type Metrics struct {
	scheduledScaleBounds m.GaugeVec
	scaleAdjustments     m.GaugeVec
}

// NewMetrics returns the metrics collection needed for the SAA.
func NewMetrics() Metrics {

	bound := []string{"bound"}
	scheduledScaleBounds := m.NewWrappedGaugeVec(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "cap",
		Name:      "scheduled_scale_bounds",
		Help:      "Shows the min and max scale value of the currently active scale schedule. In case no schedule is active both values are 0.",
	}, bound)

	aType := []string{"type"}
	scaleAdjustments := m.NewWrappedGaugeVec(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "cap",
		Name:      "adjusted_vs_planned_scale",
		Help:      "Shows the scale that was planned by the CapacityPlanner based on the current scale and the scaleFactor and the adjusted scale based on the currently active scale schedule. The value for the adjusted and the initially planned scale differs only in case the planned scale would violate a currently active scale schedule.",
	}, aType)

	return Metrics{
		scheduledScaleBounds: scheduledScaleBounds,
		scaleAdjustments:     scaleAdjustments,
	}
}
