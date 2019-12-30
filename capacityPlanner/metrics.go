package capacityPlanner

import (
	"github.com/prometheus/client_golang/prometheus"
	m "github.com/thomasobenaus/sokar/metrics"
)

// Metrics represents the collection of metrics internally set by CapacityPlanner
type Metrics struct {
	scaleAdjustments m.GaugeVec
}

// NewMetrics returns the metrics collection needed for the SAA.
func NewMetrics() Metrics {

	aType := []string{"type"}
	scaleAdjustments := m.NewWrappedGaugeVec(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "cap",
		Name:      "adjusted_vs_planned_scale",
		Help:      "Shows the scale that was planned by the CapacityPlanner based on the current scale and the scaleFactor and the adjusted scale based on the currently active scale schedule. The value for the adjusted and the initially planned scale differs only in case the planned scale would violate a currently active scale schedule.",
	}, aType)

	return Metrics{
		scaleAdjustments: scaleAdjustments,
	}
}
