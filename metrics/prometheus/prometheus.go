package prometheus

import (
	"github.com/go-kratos/kratos/v2/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

type counter struct {
	cv *prometheus.CounterVec
}

// NewCounter new a prometheus counter and returns Counter.
func NewCounter(cv *prometheus.CounterVec) metrics.Counter {
	return &counter{
		cv: cv,
	}
}

func (c *counter) Add(delta float64, labels ...string) {
	c.cv.WithLabelValues(labels...).Add(delta)
}

type gauge struct {
	gv *prometheus.GaugeVec
}

// NewGauge new a prometheus gauge and returns Gauge.
func NewGauge(gv *prometheus.GaugeVec) metrics.Gauge {
	return &gauge{
		gv: gv,
	}
}

func (g *gauge) Set(value float64, lvs ...string) {
	g.gv.WithLabelValues(lvs...).Set(value)
}

func (g *gauge) Add(delta float64, lvs ...string) {
	g.gv.WithLabelValues(lvs...).Add(delta)
}

type summary struct {
	sv *prometheus.SummaryVec
}

// NewSummary new a prometheus summary and returns Histogram.
func NewSummary(sv *prometheus.SummaryVec) metrics.Histogram {
	return &summary{
		sv: sv,
	}
}

func (s *summary) Observe(value float64, lvs ...string) {
	s.sv.WithLabelValues(lvs...).Observe(value)
}

type histogram struct {
	hv *prometheus.HistogramVec
}

// NewHistogram new a prometheus histogram and returns Histogram.
func NewHistogram(hv *prometheus.HistogramVec) metrics.Histogram {
	return &histogram{
		hv: hv,
	}
}

func (h *histogram) Observe(value float64, lvs ...string) {
	h.hv.WithLabelValues(lvs...).Observe(value)
}
