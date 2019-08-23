package metric

import (
	"errors"
	"fmt"
)

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

// VectorOpts contains the common arguments for creating vec Metric..
type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

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
func NewBusinessMetricCount(name string, labels ...string) CounterVec {
	if name == "" || len(labels) == 0 {
		panic(errors.New("stat:metric business count metric name should not be empty or labels length should be greater than zero"))
	}
	return NewCounterVec(&CounterVecOpts{
		Namespace: _businessNamespace,
		Subsystem: _businessSubsystemCount,
		Name:      name,
		Labels:    labels,
		Help:      fmt.Sprintf("business metric count %s", name),
	})
}

// NewBusinessMetricGauge business Metric gauge vec.
// name or labels should not be empty.
func NewBusinessMetricGauge(name string, labels ...string) GaugeVec {
	if name == "" || len(labels) == 0 {
		panic(errors.New("stat:metric business gauge metric name should not be empty or labels length should be greater than zero"))
	}
	return NewGaugeVec(&GaugeVecOpts{
		Namespace: _businessNamespace,
		Subsystem: _businessSubSystemGauge,
		Name:      name,
		Labels:    labels,
		Help:      fmt.Sprintf("business metric gauge %s", name),
	})
}

// NewBusinessMetricHistogram business Metric histogram vec.
// name or labels should not be empty.
func NewBusinessMetricHistogram(name string, buckets []float64, labels ...string) HistogramVec {
	if name == "" || len(labels) == 0 {
		panic(errors.New("stat:metric business histogram metric name should not be empty or labels length should be greater than zero"))
	}
	if len(buckets) == 0 {
		buckets = _defaultBuckets
	}
	return NewHistogramVec(&HistogramVecOpts{
		Namespace: _businessNamespace,
		Subsystem: _businessSubSystemHistogram,
		Name:      name,
		Labels:    labels,
		Buckets:   buckets,
		Help:      fmt.Sprintf("business metric histogram %s", name),
	})
}
