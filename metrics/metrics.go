package metrics

type Counter interface {
	Inc()
	Add(float64)
}

type Gauge interface {
	Set(float64)
	Add(float64)
}

type GaugeVec interface {
	WithLabelValues(lvs ...string) Gauge
}
