package metrics

type Counter interface {
	Inc()
}

type Gauge interface {
	Set(float64)
}
