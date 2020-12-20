package prometheus

import (
	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// counterVec counter vec.
type promCounterVec struct {
	counter *prometheus.CounterVec
}

// NewCounterVec .
func NewCounterVec(cfg *metrics.CounterVecOpts) metrics.CounterVec {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)
	prometheus.MustRegister(vec)
	return &promCounterVec{
		counter: vec,
	}
}

// Inc Inc increments the counter by 1. Use Add to increment it by arbitrary.
func (counter *promCounterVec) Inc(labels ...string) {
	counter.counter.WithLabelValues(labels...).Inc()
}

// Add Inc increments the counter by 1. Use Add to increment it by arbitrary.
func (counter *promCounterVec) Add(v float64, labels ...string) {
	counter.counter.WithLabelValues(labels...).Add(v)
}
