package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-kratos/kratos/v2/metrics"
)

var _ metrics.Counter = (*counter)(nil)

type counter struct {
	cv  *prometheus.CounterVec
	lvs []string
}

// NewCounter new a prometheus counter and returns Counter.
func NewCounter(cv *prometheus.CounterVec) metrics.Counter {
	return &counter{
		cv: cv,
	}
}

func (c *counter) With(lvs ...string) metrics.Counter {
	return &counter{
		cv:  c.cv,
		lvs: lvs,
	}
}

func (c *counter) Inc() {
	c.cv.WithLabelValues(c.lvs...).Inc()
}

func (c *counter) Add(delta float64) {
	c.cv.WithLabelValues(c.lvs...).Add(delta)
}
