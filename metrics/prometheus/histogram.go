package prometheus

import (
	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Histogram prom histogram collection.
type promHistogramVec struct {
	histogram *prometheus.HistogramVec
}

// NewHistogramVec new a histogram vec.
func NewHistogramVec(cfg *metrics.HistogramVecOpts) metrics.HistogramVec {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
			Buckets:   cfg.Buckets,
		}, cfg.Labels)
	prometheus.MustRegister(vec)
	return &promHistogramVec{
		histogram: vec,
	}
}

// Timing adds a single observation to the histogram.
func (histogram *promHistogramVec) Observe(v int64, labels ...string) {
	histogram.histogram.WithLabelValues(labels...).Observe(float64(v))
}
