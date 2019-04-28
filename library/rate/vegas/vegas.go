package vegas

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/rate"
)

const (
	_minWindowTime = int64(time.Millisecond * 500)
	_maxWindowTime = int64(time.Millisecond * 2000)

	_minLimit = 8
	_maxLimit = 2048
)

// Stat is the Statistics of vegas.
type Stat struct {
	Limit    int64
	InFlight int64
	MinRTT   time.Duration
	LastRTT  time.Duration
}

// New new a rate vegas.
func New() *Vegas {
	v := &Vegas{
		probes: 100,
		limit:  _minLimit,
	}
	v.sample.Store(&sample{})
	return v
}

// Vegas tcp vegas.
type Vegas struct {
	limit      int64
	inFlight   int64
	updateTime int64
	minRTT     int64

	sample atomic.Value
	mu     sync.Mutex
	probes int64
}

// Stat return the statistics of vegas.
func (v *Vegas) Stat() Stat {
	return Stat{
		Limit:    atomic.LoadInt64(&v.limit),
		InFlight: atomic.LoadInt64(&v.inFlight),
		MinRTT:   time.Duration(atomic.LoadInt64(&v.minRTT)),
		LastRTT:  time.Duration(v.sample.Load().(*sample).RTT()),
	}
}

// Acquire No matter success or not,done() must be called at last.
func (v *Vegas) Acquire() (done func(time.Time, rate.Op), success bool) {
	inFlight := atomic.AddInt64(&v.inFlight, 1)
	if inFlight <= atomic.LoadInt64(&v.limit) {
		success = true
	}

	return func(start time.Time, op rate.Op) {
		atomic.AddInt64(&v.inFlight, -1)
		if op == rate.Ignore {
			return
		}
		end := time.Now().UnixNano()
		rtt := end - start.UnixNano()

		s := v.sample.Load().(*sample)
		if op == rate.Drop {
			s.Add(rtt, inFlight, true)
		} else if op == rate.Success {
			s.Add(rtt, inFlight, false)
		}
		if end > atomic.LoadInt64(&v.updateTime) && s.Count() >= 16 {
			v.mu.Lock()
			defer v.mu.Unlock()
			if v.sample.Load().(*sample) != s {
				return
			}
			v.sample.Store(&sample{})

			lastRTT := s.RTT()
			if lastRTT <= 0 {
				return
			}
			updateTime := end + lastRTT*5
			if lastRTT*5 < _minWindowTime {
				updateTime = end + _minWindowTime
			} else if lastRTT*5 > _maxWindowTime {
				updateTime = end + _maxWindowTime
			}
			atomic.StoreInt64(&v.updateTime, updateTime)
			limit := atomic.LoadInt64(&v.limit)
			queue := float64(limit) * (1 - float64(v.minRTT)/float64(lastRTT))
			v.probes--
			if v.probes <= 0 {
				maxFlight := s.MaxInFlight()
				if maxFlight*2 < v.limit || maxFlight <= _minLimit {
					v.probes = 3*limit + rand.Int63n(3*limit)
					v.minRTT = lastRTT
				}
			}
			if v.minRTT == 0 || lastRTT < v.minRTT {
				v.minRTT = lastRTT
			}
			var newLimit float64
			threshold := math.Sqrt(float64(limit)) / 2
			if s.Drop() {
				newLimit = float64(limit) - threshold
			} else if s.MaxInFlight()*2 < v.limit {
				return
			} else {
				if queue < threshold {
					newLimit = float64(limit) + 6*threshold
				} else if queue < 2*threshold {
					newLimit = float64(limit) + 3*threshold
				} else if queue < 3*threshold {
					newLimit = float64(limit) + threshold
				} else if queue > 6*threshold {
					newLimit = float64(limit) - threshold
				} else {
					return
				}
			}
			newLimit = math.Max(_minLimit, math.Min(_maxLimit, newLimit))
			atomic.StoreInt64(&v.limit, int64(newLimit))
		}
	}, success
}
