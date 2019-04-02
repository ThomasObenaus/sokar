package scaler

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	m "github.com/thomasobenaus/sokar/metrics"
)

// Metrics represents the collection of metrics internally set by scaler.
type Metrics struct {
	scalingPolicyViolated        m.CounterVec
	scalingTicketCount           m.CounterVec
	scaleResultCounter           m.CounterVec
	scalingDurationSeconds       m.Histogram
	plannedButSkippedScalingOpen m.GaugeVec
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

	resultType := []string{"result"}
	scaleResultCounter := m.NewWrappedCounterVec(prometheus.CounterOpts{
		Namespace: "sokar",
		Subsystem: "sca",
		Name:      "scale_result_counter",
		Help:      "Counts the result types of a scaling action (success, failed, ignored).",
	}, resultType)

	scalingDurationSeconds := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "sokar",
		Subsystem: "sca",
		Name:      "scaling_duration_seconds",
		Help:      "Holds the duration of the scaling actions so far. This is the time it took to apply a scaling (execute the deployment).",
		Buckets:   []float64{0.2, 0.5, 1, 2, 5, 8, 15, 20, 25, 30, 40, 50, 75, 100},
	})

	direction := []string{"direction"}
	plannedButSkippedScalingOpen := m.NewWrappedGaugeVec(prometheus.GaugeOpts{
		Namespace: "sokar",
		Subsystem: "sca",
		Name:      "planned_but_skipped_scaling_open",
		Help:      "Is a helper metric which is only used in dry run mode. It is set to 1 in case there was a automatic scaling planned but not exectued due to dry-run mode. It is reset to 0 if then a scaling was applied.",
	}, direction)

	return Metrics{
		scalingPolicyViolated:        scalingPolicyViolated,
		scalingTicketCount:           scalingTicketCount,
		scaleResultCounter:           scaleResultCounter,
		scalingDurationSeconds:       scalingDurationSeconds,
		plannedButSkippedScalingOpen: plannedButSkippedScalingOpen,
	}
}
