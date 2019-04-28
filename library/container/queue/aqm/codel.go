package aqm

import (
	"context"
	"math"
	"sync"
	"time"

	"go-common/library/ecode"
)

// Config codel config.
type Config struct {
	Target   int64 // target queue delay (default 20 ms).
	Internal int64 // sliding minimum time window width (default 500 ms)
}

// Stat is the Statistics of codel.
type Stat struct {
	Dropping bool
	FaTime   int64
	DropNext int64
	Packets  int
}

type packet struct {
	ch chan bool
	ts int64
}

var defaultConf = &Config{
	Target:   50,
	Internal: 500,
}

// Queue queue is CoDel req buffer queue.
type Queue struct {
	pool    sync.Pool
	packets chan packet

	mux      sync.RWMutex
	conf     *Config
	count    int64
	dropping bool  // 	Equal to 1 if in drop state
	faTime   int64 // Time when we'll declare we're above target (0 if below)
	dropNext int64 // Packets dropped since going into drop state
}

// Default new a default codel queue.
func Default() *Queue {
	return New(defaultConf)
}

// New new codel queue.
func New(conf *Config) *Queue {
	if conf == nil {
		conf = defaultConf
	}
	q := &Queue{
		packets: make(chan packet, 2048),
		conf:    conf,
	}
	q.pool.New = func() interface{} {
		return make(chan bool)
	}
	return q
}

// Reload set queue config.
func (q *Queue) Reload(c *Config) {
	if c == nil || c.Internal <= 0 || c.Target <= 0 {
		return
	}
	// TODO codel queue size
	q.mux.Lock()
	q.conf = c
	q.mux.Unlock()
}

// Stat return the statistics of codel
func (q *Queue) Stat() Stat {
	q.mux.Lock()
	defer q.mux.Unlock()
	return Stat{
		Dropping: q.dropping,
		FaTime:   q.faTime,
		DropNext: q.dropNext,
		Packets:  len(q.packets),
	}
}

// Push req into CoDel request buffer queue.
// if return error is nil,the caller must call q.Done() after finish request handling
func (q *Queue) Push(ctx context.Context) (err error) {
	r := packet{
		ch: q.pool.Get().(chan bool),
		ts: time.Now().UnixNano() / int64(time.Millisecond),
	}
	select {
	case q.packets <- r:
	default:
		err = ecode.LimitExceed
		q.pool.Put(r.ch)
	}
	if err == nil {
		select {
		case drop := <-r.ch:
			if drop {
				err = ecode.LimitExceed
			}
			q.pool.Put(r.ch)
		case <-ctx.Done():
			err = ecode.Deadline
		}
	}
	return
}

// Pop req from CoDel request buffer queue.
func (q *Queue) Pop() {
	for {
		select {
		case p := <-q.packets:
			drop := q.judge(p)
			select {
			case p.ch <- drop:
				if !drop {
					return
				}
			default:
				q.pool.Put(p.ch)
			}
		default:
			return
		}
	}
}

func (q *Queue) controlLaw(now int64) int64 {
	q.dropNext = now + int64(float64(q.conf.Internal)/math.Sqrt(float64(q.count)))
	return q.dropNext
}

// judge decide if the packet should drop or not.
func (q *Queue) judge(p packet) (drop bool) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	sojurn := now - p.ts
	q.mux.Lock()
	defer q.mux.Unlock()
	if sojurn < q.conf.Target {
		q.faTime = 0
	} else if q.faTime == 0 {
		q.faTime = now + q.conf.Internal
	} else if now >= q.faTime {
		drop = true
	}
	if q.dropping {
		if !drop {
			// sojourn time below target - leave dropping state
			q.dropping = false
		} else if now > q.dropNext {
			q.count++
			q.dropNext = q.controlLaw(q.dropNext)
			drop = true
			return
		}
	} else if drop && (now-q.dropNext < q.conf.Internal || now-q.faTime >= q.conf.Internal) {
		q.dropping = true
		// If we're in a drop cycle, the drop rate that controlled the queue
		// on the last cycle is a good starting point to control it now.
		if now-q.dropNext < q.conf.Internal {
			if q.count > 2 {
				q.count = q.count - 2
			} else {
				q.count = 1
			}
		} else {
			q.count = 1
		}
		q.dropNext = q.controlLaw(now)
		drop = true
		return
	}
	return
}
