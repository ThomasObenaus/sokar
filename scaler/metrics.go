package scaler

import (
	"github.com/prometheus/client_golang/prometheus"
	m "github.com/thomasobenaus/sokar/metrics"
)

// Metrics represents the collection of metrics internally set by scaler.
type Metrics struct {
	scalingPolicyViolated m.CounterVec
	scalingTicketCount    m.CounterVec
}

// NewMetrics returns the metrics collection needed for the SAA.
func NewMetrics() Metrics {

	thresholds := []string{"threshold"}
	scalingPolicyViolated := m.NewWrappedCounterVec(prometheus.CounterOpts{
		Namespace: "sokar",
		Subsystem: "sca",
		Name:      "scaling_policy_violated",
		Help:      "Counts the number of occurrences the planning of sokar would have violated the scaling policy of the job (upper or lower threshold).",
	}, thresholds)

	ticketAction := []string{"action"}
	scalingTicketCount := m.NewWrappedCounterVec(prometheus.CounterOpts{
		Namespace: "sokar",
		Subsystem: "sca",
		Name:      "scaling_ticket_counter",
		Help:      "Counts the number of added, rejected and applied scaling tickets.",
	}, ticketAction)

	return Metrics{
		scalingPolicyViolated: scalingPolicyViolated,
		scalingTicketCount:    scalingTicketCount,
	}
}
