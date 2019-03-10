package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// WrappedGaugeVec wraps a prometheus GaugeVec
type WrappedGaugeVec struct {
	prom *prometheus.GaugeVec
}

// WithLabelValues implements the WithLabelValues to meet the GaugeVec interface
func (wG *WrappedGaugeVec) WithLabelValues(lvs ...string) Gauge {
	return wG.prom.WithLabelValues(lvs...)
}

// NewWrappedGaugeVec creates a prometheus GaugeVec that is wrapped
func NewWrappedGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *WrappedGaugeVec {
	return &WrappedGaugeVec{
		prom: promauto.NewGaugeVec(opts, labelNames),
	}
}
