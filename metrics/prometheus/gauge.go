package prometheus

import (
	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// gaugeVec gauge vec.
type promGaugeVec struct {
	gauge *prometheus.GaugeVec
}

// NewGaugeVec .
func NewGaugeVec(cfg *metrics.GaugeVecOpts) metrics.GaugeVec {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)
	prometheus.MustRegister(vec)
	return &promGaugeVec{
		gauge: vec,
	}
}

// Inc Inc increments the counter by 1. Use Add to increment it by arbitrary.
func (gauge *promGaugeVec) Inc(labels ...string) {
	gauge.gauge.WithLabelValues(labels...).Inc()
}

// Add Inc increments the counter by 1. Use Add to increment it by arbitrary.
func (gauge *promGaugeVec) Add(v float64, labels ...string) {
	gauge.gauge.WithLabelValues(labels...).Add(v)
}

// Set set the given value to the collection.
func (gauge *promGaugeVec) Set(v float64, labels ...string) {
	gauge.gauge.WithLabelValues(labels...).Set(v)
}
