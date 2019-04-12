package breaker

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/bilibili/kratos/pkg/ecode"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/stat/metric"
)

// sreBreaker is a sre CircuitBreaker pattern.
type sreBreaker struct {
	stat metric.RollingCounter

	k       float64
	request int64

	state int32
	r     *rand.Rand
}

func newSRE(c *Config) Breaker {
	counterOpts := metric.RollingCounterOpts{
		Size:           c.Bucket,
		BucketDuration: time.Duration(int64(c.Window) / int64(c.Bucket)),
	}
	stat := metric.NewRollingCounter(counterOpts)
	return &sreBreaker{
		stat: stat,
		r:    rand.New(rand.NewSource(time.Now().UnixNano())),

		request: c.Request,
		k:       c.K,
		state:   StateClosed,
	}
}

func (b *sreBreaker) summary() (success int64, total int64) {
	b.stat.Reduce(func(iterator metric.Iterator) float64 {
		for iterator.Next() {
			bucket := iterator.Bucket()
			total += bucket.Count
			for _, p := range bucket.Points {
				success += int64(p)
			}
		}
		return 0
	})
	return
}

func (b *sreBreaker) Allow() error {
	success, total := b.summary()
	k := b.k * float64(success)
	if log.V(5) {
		log.Info("breaker: request: %d, succee: %d, fail: %d", total, success, total-success)
	}
	// check overflow requests = K * success
	if total < b.request || float64(total) < k {
		if atomic.LoadInt32(&b.state) == StateOpen {
			atomic.CompareAndSwapInt32(&b.state, StateOpen, StateClosed)
		}
		return nil
	}
	if atomic.LoadInt32(&b.state) == StateClosed {
		atomic.CompareAndSwapInt32(&b.state, StateClosed, StateOpen)
	}
	dr := math.Max(0, (float64(total)-k)/float64(total+1))
	rr := b.r.Float64()
	if log.V(5) {
		log.Info("breaker: drop ratio: %f, real rand: %f, drop: %v", dr, rr, dr > rr)
	}
	if dr <= rr {
		return nil
	}
	return ecode.ServiceUnavailable
}

func (b *sreBreaker) MarkSuccess() {
	b.stat.Add(1)
}

func (b *sreBreaker) MarkFailed() {
	// NOTE: when client reject requets locally, continue add counter let the
	// drop ratio higher.
	b.stat.Add(0)
}
