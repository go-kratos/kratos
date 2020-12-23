package metrics

// HistogramVecOpts is histogram vector opts.
type HistogramVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	Buckets   []float64
}

// HistogramVec gauge vec.
type HistogramVec interface {
	// Observe adds a single observation to the histogram.
	Observe(v int64, labels ...string)
}
