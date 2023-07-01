package metrics

// Counter is metrics counter.
// Deprecated: use otel metrcis instead.
type Counter interface {
	With(lvs ...string) Counter
	Inc()
	Add(delta float64)
}

// Gauge is metrics gauge.
// Deprecated: use otel metrcis instead.
type Gauge interface {
	With(lvs ...string) Gauge
	Set(value float64)
	Add(delta float64)
	Sub(delta float64)
}

// Observer is metrics observer.
// Deprecated: use otel metrcis instead.
type Observer interface {
	With(lvs ...string) Observer
	Observe(float64)
}
