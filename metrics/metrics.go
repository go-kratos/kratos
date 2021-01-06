package metrics

// Counter is metrics counter.
type Counter interface {
	Add(delta float64, lvs ...string)
}

// Gauge is metrics gauge.
type Gauge interface {
	Set(value float64, lvs ...string)
	Add(delta float64, lvs ...string)
}

// Histogram is metrics histogram.
type Histogram interface {
	Observe(value float64, lvs ...string)
}
