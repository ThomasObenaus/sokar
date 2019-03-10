package metrics

type Counter interface {
	Inc()
}

type Gauge interface {
	Set(float64)
}

type GaugeVec interface {
	WithLabelValues(lvs ...string) Gauge
}
