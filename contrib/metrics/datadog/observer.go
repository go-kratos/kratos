package metrics

import (
	"time"

	"github.com/go-kratos/kratos/v2/metrics"
)

var _ metrics.Observer = (*timing)(nil)

type timing struct {
	opts options
	name string
	lvs  []string
}

// NewTiming new a DataDog timer and returns Observer.
func NewTiming(name string, opts ...Option) metrics.Observer {
	options := options{
		sampleRate: 1,
		client:     defaultClient,
	}
	for _, o := range opts {
		o(&options)
	}
	return &timing{
		name: name,
		opts: options,
	}
}

// With is applied in kratos/middleware/metrics/metrics.go (method,path)
func (d *timing) With(values ...string) metrics.Observer {
	return &timing{
		name: d.name,
		opts: d.opts,
		lvs:  withValues(d.opts.labels, values),
	}
}

func (d *timing) Observe(value float64) {
	d.opts.client.TimeInMilliseconds(d.name, value*float64(time.Second/time.Millisecond), d.lvs, d.opts.sampleRate)
}
