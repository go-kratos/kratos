package metrics

import (
	"github.com/go-kratos/kratos/v2/metrics"
)

var _ metrics.Counter = (*counter)(nil)

type counter struct {
	opts options
	name string
	lvs  []string
}

// NewCounter new a DataDog counter and returns Counter.
func NewCounter(name string, opts ...Option) metrics.Counter {
	options := options{
		sampleRate: 1,
		client:     defaultClient,
	}
	for _, o := range opts {
		o(&options)
	}
	return &counter{
		name: name,
		opts: options,
	}
}

// With is applied in kratos/middleware/metrics/metrics.go (method,path,code)
func (d *counter) With(values ...string) metrics.Counter {
	return &counter{
		name: d.name,
		opts: d.opts,
		lvs:  withValues(d.opts.labels, values),
	}
}

func (d *counter) Inc() {
	d.opts.client.Incr(d.name, d.lvs, d.opts.sampleRate)
}

func (d *counter) Add(delta float64) {
	d.opts.client.Count(d.name, int64(delta), d.lvs, d.opts.sampleRate)
}
