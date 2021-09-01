package metrics

import (
	"github.com/go-kratos/kratos/v2/metrics"
)

var _ metrics.Gauge = (*gauge)(nil)

type gauge struct {
	opts options
	name string
	lvs  []string
}

// NewGauge new a DataDog gauge and returns Gauge.
func NewGauge(name string, opts ...Option) metrics.Gauge {
	options := options{
		sampleRate: 1,
		client:     defaultClient,
	}
	for _, o := range opts {
		o(&options)
	}
	return &gauge{
		name: name,
		opts: options,
	}
}

func (d *gauge) With(values ...string) metrics.Gauge {
	return &gauge{
		opts: d.opts,
		name: d.name,
		lvs:  withValues(d.opts.labels, values),
	}
}

func (d *gauge) Set(value float64) {
	d.opts.client.Gauge(d.name, value, d.lvs, d.opts.sampleRate)
}

func (d *gauge) Add(delta float64) {
}

func (d *gauge) Sub(delta float64) {
}
