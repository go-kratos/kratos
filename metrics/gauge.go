package metrics

// Gauge stores a numerical value that can be add arbitrarily.
type Gauge interface {
	Metric
	// Sets sets the value to the given number.
	Set(int64)
}

// GaugeVec gauge vec.
type GaugeVec interface {
	// Set sets the Gauge to an arbitrary value.
	Set(v float64, labels ...string)
	// Inc increments the Gauge by 1. Use Add to increment it by arbitrary
	// values.
	Inc(labels ...string)
	// Add adds the given value to the Gauge. (The value can be negative,
	// resulting in a decrease of the Gauge.)
	Add(v float64, labels ...string)
}
