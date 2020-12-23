package counter

import (
	"fmt"
	"sync/atomic"

	"github.com/go-kratos/kratos/v2/metrics"
)

var _ metrics.Metric = &counter{}

// CounterOpts is an alias of Opts.
type CounterOpts metrics.Opts

type counter struct {
	val int64
}

// NewCounter creates a new Counter based on the CounterOpts.
func NewCounter(opts CounterOpts) metrics.Counter {
	return &counter{}
}

func (c *counter) Add(val int64) {
	if val < 0 {
		panic(fmt.Errorf("stat/metric: cannot decrease in negative value. val: %d", val))
	}
	atomic.AddInt64(&c.val, val)
}

func (c *counter) Value() int64 {
	return atomic.LoadInt64(&c.val)
}
