package vegas

import (
	"sync/atomic"
)

type sample struct {
	count       int64
	maxInFlight int64
	drop        int64
	// nanoseconds
	totalRTT int64
}

func (s *sample) Add(rtt int64, inFlight int64, drop bool) {
	if drop {
		atomic.StoreInt64(&s.drop, 1)
	}
	for max := atomic.LoadInt64(&s.maxInFlight); max < inFlight; max = atomic.LoadInt64(&s.maxInFlight) {
		if atomic.CompareAndSwapInt64(&s.maxInFlight, max, inFlight) {
			break
		}
	}
	atomic.AddInt64(&s.totalRTT, rtt)
	atomic.AddInt64(&s.count, 1)
}

func (s *sample) RTT() int64 {
	count := atomic.LoadInt64(&s.count)
	if count == 0 {
		return 0
	}
	return atomic.LoadInt64(&s.totalRTT) / count
}

func (s *sample) MaxInFlight() int64 {
	return atomic.LoadInt64(&s.maxInFlight)
}

func (s *sample) Count() int64 {
	return atomic.LoadInt64(&s.count)
}

func (s *sample) Drop() bool {
	return atomic.LoadInt64(&s.drop) == 1
}
