package summary

import (
	"sync"
	"time"
)

type bucket struct {
	val   int64
	count int64
	next  *bucket
}

func (b *bucket) Add(val int64) {
	b.val += val
	b.count++
}

func (b *bucket) Value() (int64, int64) {
	return b.val, b.count
}

func (b *bucket) Reset() {
	b.val = 0
	b.count = 0
}

// Summary is a summary interface.
type Summary interface {
	Add(int64)
	Reset()
	Value() (val int64, cnt int64)
}

type summary struct {
	mu         sync.RWMutex
	buckets    []bucket
	bucketTime int64
	lastAccess int64
	cur        *bucket
}

// New new a summary.
//
// use RollingCounter creates a new window. windowTime is the time covering the entire
// window. windowBuckets is the number of buckets the window is divided into.
// An example: a 10 second window with 10 buckets will have 10 buckets covering
// 1 second each.
func New(window time.Duration, winBucket int) Summary {
	buckets := make([]bucket, winBucket)
	bucket := &buckets[0]
	for i := 1; i < winBucket; i++ {
		bucket.next = &buckets[i]
		bucket = bucket.next
	}
	bucket.next = &buckets[0]
	bucketTime := time.Duration(window.Nanoseconds() / int64(winBucket))
	return &summary{
		cur:        &buckets[0],
		buckets:    buckets,
		bucketTime: int64(bucketTime),
		lastAccess: time.Now().UnixNano(),
	}
}

// Add increments the summary by value.
func (s *summary) Add(val int64) {
	s.mu.Lock()
	s.lastBucket().Add(val)
	s.mu.Unlock()
}

// Value get the summary value and count.
func (s *summary) Value() (val int64, cnt int64) {
	now := time.Now().UnixNano()
	s.mu.RLock()
	b := s.cur
	i := s.elapsed(now)
	for j := 0; j < len(s.buckets); j++ {
		// skip all future reset bucket.
		if i > 0 {
			i--
		} else {
			v, c := b.Value()
			val += v
			cnt += c
		}
		b = b.next
	}
	s.mu.RUnlock()
	return
}

//  Reset reset the counter.
func (s *summary) Reset() {
	s.mu.Lock()
	for i := range s.buckets {
		s.buckets[i].Reset()
	}
	s.mu.Unlock()
}

func (s *summary) elapsed(now int64) (i int) {
	var e int64
	if e = now - s.lastAccess; e <= s.bucketTime {
		return
	}
	if i = int(e / s.bucketTime); i > len(s.buckets) {
		i = len(s.buckets)
	}
	return
}

func (s *summary) lastBucket() (b *bucket) {
	now := time.Now().UnixNano()
	b = s.cur
	// reset the buckets between now and number of buckets ago. If
	// that is more that the existing buckets, reset all.
	if i := s.elapsed(now); i > 0 {
		s.lastAccess = now
		for ; i > 0; i-- {
			// replace the next used bucket.
			b = b.next
			b.Reset()
		}
	}
	s.cur = b
	return
}
