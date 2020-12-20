package prometheus

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/pkg/errors"
)

const (
	_businessNamespace          = "business"
	_businessSubsystemCount     = "count"
	_businessSubSystemGauge     = "gauge"
	_businessSubSystemHistogram = "histogram"
)

var (
	_defaultBuckets = []float64{5, 10, 25, 50, 100, 250, 500}
)

// NewBusinessMetricCount business Metric count vec.
// name or labels should not be empty.
func NewBusinessMetricCount(name string, labels ...string) metrics.CounterVec {
	if name == "" || len(labels) == 0 {
		panic(errors.New("stat:metric business count metric name should not be empty or labels length should be greater than zero"))
	}
	return NewCounterVec(&metrics.CounterVecOpts{
		Namespace: _businessNamespace,
		Subsystem: _businessSubsystemCount,
		Name:      name,
		Labels:    labels,
		Help:      fmt.Sprintf("business metric count %s", name),
	})
}

// NewBusinessMetricGauge business Metric gauge vec.
// name or labels should not be empty.
func NewBusinessMetricGauge(name string, labels ...string) metrics.GaugeVec {
	if name == "" || len(labels) == 0 {
		panic(errors.New("stat:metric business gauge metric name should not be empty or labels length should be greater than zero"))
	}
	return NewGaugeVec(&metrics.GaugeVecOpts{
		Namespace: _businessNamespace,
		Subsystem: _businessSubSystemGauge,
		Name:      name,
		Labels:    labels,
		Help:      fmt.Sprintf("business metric gauge %s", name),
	})
}

// NewBusinessMetricHistogram business Metric histogram vec.
// name or labels should not be empty.
func NewBusinessMetricHistogram(name string, buckets []float64, labels ...string) metrics.HistogramVec {
	if name == "" || len(labels) == 0 {
		panic(errors.New("stat:metric business histogram metric name should not be empty or labels length should be greater than zero"))
	}
	if len(buckets) == 0 {
		buckets = _defaultBuckets
	}
	return NewHistogramVec(&metrics.HistogramVecOpts{
		Namespace: _businessNamespace,
		Subsystem: _businessSubSystemHistogram,
		Name:      name,
		Labels:    labels,
		Buckets:   buckets,
		Help:      fmt.Sprintf("business metric histogram %s", name),
	})
}
