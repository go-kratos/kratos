package metric

// Opts contains the common arguments for creating Metric.
type Opts struct {
}

// Metric is a sample interface.
// Implementations of Metrics in metric package are Counter, Gauge,
// PointGauge, RollingCounter and RollingGauge.
type Metric interface {
	// Add adds the given value to the counter.
	Add(int64)
	// Value gets the current value.
	// If the metric's type is PointGauge, RollingCounter, RollingGauge,
	// it returns the sum value within the window.
	Value() int64
}

// Aggregation contains some common aggregation function.
// Each aggregation can compute summary statistics of window.
type Aggregation interface {
	// Min finds the min value within the window.
	Min() float64
	// Max finds the max value within the window.
	Max() float64
	// Avg computes average value within the window.
	Avg() float64
	// Sum computes sum value within the window.
	Sum() float64
}
