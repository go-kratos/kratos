package limit

import (
	"context"
	"time"

	"go-common/library/container/queue/aqm"
	"go-common/library/log"
	"go-common/library/rate"
	"go-common/library/rate/vegas"
)

var _ rate.Limiter = &Limiter{}

// New returns a new Limiter that allows events up to adaptive rtt.
func New(c *aqm.Config) *Limiter {
	l := &Limiter{
		rate:  vegas.New(),
		queue: aqm.New(c),
	}
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for {
			<-ticker.C
			v := l.rate.Stat()
			q := l.queue.Stat()
			log.Info("rate/limit: limit(%d) inFlight(%d) minRtt(%v) rtt(%v) codel packets(%d)", v.Limit, v.InFlight, v.MinRTT, v.LastRTT, q.Packets)
		}
	}()
	return l
}

// Limiter use tcp vegas + codel for adaptive limit.
type Limiter struct {
	rate  *vegas.Vegas
	queue *aqm.Queue
}

// Allow immplemnet rate.Limiter.
// if error is returned,no need to call done()
func (l *Limiter) Allow(ctx context.Context) (func(rate.Op), error) {
	var (
		done func(time.Time, rate.Op)
		err  error
		ok   bool
	)
	if done, ok = l.rate.Acquire(); !ok {
		// NOTE exceed max inflight, use queue
		if err = l.queue.Push(ctx); err != nil {
			done(time.Time{}, rate.Ignore)
			return func(rate.Op) {}, err
		}
	}
	start := time.Now()
	return func(op rate.Op) {
		done(start, op)
		l.queue.Pop()
	}, nil
}
