package circuitbreaker

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// StateOpen rejects requests according to the calculated drop ratio.
	StateOpen int32 = iota
	// StateClosed allows requests while the rolling failure ratio is healthy.
	StateClosed
)

// Option configures the SRE circuit breaker.
type Option func(*options)

type options struct {
	failureRatio float64
	request      int64
	bucket       int
	window       time.Duration
}

// WithFailureRatio sets the failure ratio threshold that starts rejection.
func WithFailureRatio(ratio float64) Option {
	return func(o *options) {
		o.failureRatio = ratio
	}
}

// WithRequest sets the minimum request count before rejection starts.
func WithRequest(r int64) Option {
	return func(o *options) {
		o.request = r
	}
}

// WithWindow sets the rolling statistical window.
func WithWindow(d time.Duration) Option {
	return func(o *options) {
		o.window = d
	}
}

// WithBucket sets the number of buckets in the rolling window.
func WithBucket(b int) Option {
	return func(o *options) {
		o.bucket = b
	}
}

// Breaker is an SRE-style circuit breaker.
type Breaker struct {
	stat *rollingCounter

	random func() float64

	k       float64
	request int64
	state   int32
}

// NewBreaker returns an SRE circuit breaker.
func NewBreaker(opts ...Option) CircuitBreaker {
	opt := options{
		failureRatio: 0.5,
		request:      20,
		bucket:       10,
		window:       3 * time.Second,
	}
	for _, o := range opts {
		o(&opt)
	}
	if opt.failureRatio < 0 || opt.failureRatio >= 1 {
		opt.failureRatio = 0.5
	}
	if opt.request < 1 {
		opt.request = 1
	}
	if opt.bucket < 1 {
		opt.bucket = 1
	}
	if opt.window <= 0 {
		opt.window = 3 * time.Second
	}
	bucketDuration := opt.window / time.Duration(opt.bucket)
	if bucketDuration <= 0 {
		bucketDuration = opt.window
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randMu sync.Mutex
	return &Breaker{
		stat:    newRollingCounter(opt.bucket, bucketDuration),
		random:  func() float64 { randMu.Lock(); defer randMu.Unlock(); return rnd.Float64() },
		request: opt.request,
		k:       1 / (1 - opt.failureRatio),
		state:   StateClosed,
	}
}

// Allow reports whether the request can pass the breaker.
func (b *Breaker) Allow() error {
	successes, total := b.stat.summary()
	requests := b.k * float64(successes)
	if total < b.request || float64(total) < requests {
		atomic.CompareAndSwapInt32(&b.state, StateOpen, StateClosed)
		return nil
	}
	atomic.CompareAndSwapInt32(&b.state, StateClosed, StateOpen)
	dropRatio := math.Max(0, (float64(total)-requests)/float64(total+1))
	if b.random() < dropRatio {
		return ErrNotAllowed
	}
	return nil
}

// MarkSuccess records a successful request.
func (b *Breaker) MarkSuccess() {
	b.stat.add(1)
}

// MarkFailed records a failed request.
func (b *Breaker) MarkFailed() {
	b.stat.add(0)
}

type rollingCounter struct {
	mu             sync.Mutex
	buckets        []counterBucket
	bucketDuration time.Duration
}

type counterBucket struct {
	slot    int64
	success int64
	total   int64
}

func newRollingCounter(size int, bucketDuration time.Duration) *rollingCounter {
	return &rollingCounter{
		buckets:        make([]counterBucket, size),
		bucketDuration: bucketDuration,
	}
}

func (r *rollingCounter) add(success int64) {
	slot := r.currentSlot()
	offset := int(slot % int64(len(r.buckets)))

	r.mu.Lock()
	defer r.mu.Unlock()
	bucket := &r.buckets[offset]
	if bucket.slot != slot {
		bucket.slot = slot
		bucket.success = 0
		bucket.total = 0
	}
	bucket.success += success
	bucket.total++
}

func (r *rollingCounter) summary() (success int64, total int64) {
	slot := r.currentSlot()
	size := int64(len(r.buckets))

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, bucket := range r.buckets {
		if bucket.total == 0 || slot-bucket.slot >= size || bucket.slot > slot {
			continue
		}
		success += bucket.success
		total += bucket.total
	}
	return success, total
}

func (r *rollingCounter) currentSlot() int64 {
	return time.Now().UnixNano() / int64(r.bucketDuration)
}
