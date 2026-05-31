package ratelimit

import (
	"math"
	"runtime"
	"runtime/metrics"
	"sync"
	"sync/atomic"
	"time"
)

var (
	gCPU  int64
	decay = 0.95
)

type (
	cpuGetter func() int64

	// Option configures a BBR limiter.
	Option func(*options)
)

func init() {
	go cpuproc()
}

func cpuproc() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	sampler := newCPUSampler()
	for range ticker.C {
		curCPU := sampler.usage()
		if curCPU < 0 {
			continue
		}
		curCPU = minInt64(curCPU, 1000)
		prevCPU := atomic.LoadInt64(&gCPU)
		cpu := int64(float64(prevCPU)*decay + float64(curCPU)*(1.0-decay))
		atomic.StoreInt64(&gCPU, cpu)
	}
}

type cpuSampler struct {
	mu        sync.Mutex
	samples   []metrics.Sample
	prevUsed  float64
	prevTotal float64
	ready     bool
}

func newCPUSampler() *cpuSampler {
	return &cpuSampler{
		samples: []metrics.Sample{
			{Name: "/cpu/classes/total:cpu-seconds"},
			{Name: "/cpu/classes/idle:cpu-seconds"},
		},
	}
}

func (s *cpuSampler) usage() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	metrics.Read(s.samples)
	total := s.samples[0].Value.Float64()
	idle := s.samples[1].Value.Float64()
	used := total - idle
	if !s.ready {
		s.prevUsed = used
		s.prevTotal = total
		s.ready = true
		return 0
	}
	usedDelta := used - s.prevUsed
	totalDelta := total - s.prevTotal
	s.prevUsed = used
	s.prevTotal = total
	if totalDelta <= 0 || usedDelta < 0 {
		return 0
	}
	return int64((usedDelta / totalDelta) * 1000)
}

func minInt64(l, r int64) int64 {
	if l < r {
		return l
	}
	return r
}

// Stat contains a BBR metrics snapshot.
type Stat struct {
	CPU         int64
	InFlight    int64
	MaxInFlight int64
	MinRt       int64
	MaxPass     int64
}

type counterCache struct {
	val  int64
	time time.Time
}

type options struct {
	Window       time.Duration
	Bucket       int
	CPUThreshold int64
	CPUQuota     float64
}

// WithWindow sets the rolling window duration.
func WithWindow(d time.Duration) Option {
	return func(o *options) {
		o.Window = d
	}
}

// WithBucket sets the rolling window bucket count.
func WithBucket(b int) Option {
	return func(o *options) {
		o.Bucket = b
	}
}

// WithCPUThreshold sets the CPU threshold, scaled from 0 to 1000.
func WithCPUThreshold(threshold int64) Option {
	return func(o *options) {
		o.CPUThreshold = threshold
	}
}

// WithCPUQuota sets the real CPU quota if it differs from GOMAXPROCS.
func WithCPUQuota(quota float64) Option {
	return func(o *options) {
		o.CPUQuota = quota
	}
}

// BBR implements a BBR-like adaptive limiter.
type BBR struct {
	cpu             cpuGetter
	passStat        *rollingCounter
	rtStat          *rollingCounter
	inFlight        int64
	bucketPerSecond int64
	bucketDuration  time.Duration

	prevDropTime atomic.Value
	maxPASSCache atomic.Value
	minRtCache   atomic.Value

	opts options
}

// NewLimiter returns a BBR limiter.
func NewLimiter(opts ...Option) *BBR {
	opt := options{
		Window:       10 * time.Second,
		Bucket:       100,
		CPUThreshold: 800,
	}
	for _, o := range opts {
		o(&opt)
	}
	if opt.Window <= 0 {
		opt.Window = 10 * time.Second
	}
	if opt.Bucket < 1 {
		opt.Bucket = 1
	}

	bucketDuration := opt.Window / time.Duration(opt.Bucket)
	if bucketDuration <= 0 {
		bucketDuration = opt.Window
	}

	limiter := &BBR{
		opts:            opt,
		passStat:        newRollingCounter(opt.Bucket, bucketDuration),
		rtStat:          newRollingCounter(opt.Bucket, bucketDuration),
		bucketDuration:  bucketDuration,
		bucketPerSecond: int64(time.Second / bucketDuration),
		cpu:             func() int64 { return atomic.LoadInt64(&gCPU) },
	}
	if limiter.bucketPerSecond < 1 {
		limiter.bucketPerSecond = 1
	}
	if opt.CPUQuota != 0 {
		limiter.cpu = func() int64 {
			return int64(float64(atomic.LoadInt64(&gCPU)) * float64(runtime.GOMAXPROCS(0)) / opt.CPUQuota)
		}
	}
	return limiter
}

func (l *BBR) maxPASS() int64 {
	passCache := l.maxPASSCache.Load()
	if passCache != nil {
		ps := passCache.(*counterCache)
		if l.timespan(ps.time) < 1 {
			return ps.val
		}
	}
	rawMaxPass := int64(l.passStat.reduce(func(bucket counterBucket) float64 {
		return float64(bucket.count)
	}, math.Max, 1))
	l.maxPASSCache.Store(&counterCache{
		val:  rawMaxPass,
		time: time.Now(),
	})
	return rawMaxPass
}

func (l *BBR) timespan(lastTime time.Time) int {
	v := int(time.Since(lastTime) / l.bucketDuration)
	if v > -1 {
		return v
	}
	return l.opts.Bucket
}

func (l *BBR) minRT() int64 {
	rtCache := l.minRtCache.Load()
	if rtCache != nil {
		rc := rtCache.(*counterCache)
		if l.timespan(rc.time) < 1 {
			return rc.val
		}
	}
	rawRT := l.rtStat.reduce(func(bucket counterBucket) float64 {
		if bucket.count == 0 {
			return math.MaxFloat64
		}
		return bucket.sum / float64(bucket.count)
	}, math.Min, math.MaxFloat64)
	rawMinRT := int64(1)
	if rawRT > 0 && rawRT != math.MaxFloat64 {
		rawMinRT = int64(math.Ceil(rawRT))
	}
	l.minRtCache.Store(&counterCache{
		val:  rawMinRT,
		time: time.Now(),
	})
	return rawMinRT
}

func (l *BBR) maxInFlight() int64 {
	return int64(math.Floor(float64(l.maxPASS()*l.minRT()*l.bucketPerSecond)/1000.0) + 0.5)
}

func (l *BBR) shouldDrop() bool {
	now := time.Duration(time.Now().UnixNano())
	if l.cpu() < l.opts.CPUThreshold {
		prevDropTime, _ := l.prevDropTime.Load().(time.Duration)
		if prevDropTime == 0 {
			return false
		}
		if now-prevDropTime <= time.Second {
			inFlight := atomic.LoadInt64(&l.inFlight)
			return inFlight > 1 && inFlight > l.maxInFlight()
		}
		l.prevDropTime.Store(time.Duration(0))
		return false
	}
	inFlight := atomic.LoadInt64(&l.inFlight)
	drop := inFlight > 1 && inFlight > l.maxInFlight()
	if drop {
		prevDrop, _ := l.prevDropTime.Load().(time.Duration)
		if prevDrop != 0 {
			return true
		}
		l.prevDropTime.Store(now)
	}
	return drop
}

// Stat returns a metrics snapshot.
func (l *BBR) Stat() Stat {
	return Stat{
		CPU:         l.cpu(),
		MinRt:       l.minRT(),
		MaxPass:     l.maxPASS(),
		MaxInFlight: l.maxInFlight(),
		InFlight:    atomic.LoadInt64(&l.inFlight),
	}
}

// Allow checks whether a request is allowed.
func (l *BBR) Allow() (DoneFunc, error) {
	if l.shouldDrop() {
		return nil, ErrLimitExceed
	}
	atomic.AddInt64(&l.inFlight, 1)
	start := time.Now()
	return func(DoneInfo) {
		if rt := math.Ceil(float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond)); rt > 0 {
			l.rtStat.add(rt)
		}
		atomic.AddInt64(&l.inFlight, -1)
		l.passStat.add(1)
	}, nil
}

type rollingCounter struct {
	mu             sync.Mutex
	buckets        []counterBucket
	bucketDuration time.Duration
}

type counterBucket struct {
	slot  int64
	sum   float64
	count int64
}

func newRollingCounter(size int, bucketDuration time.Duration) *rollingCounter {
	return &rollingCounter{
		buckets:        make([]counterBucket, size),
		bucketDuration: bucketDuration,
	}
}

func (r *rollingCounter) add(value float64) {
	slot := r.currentSlot()
	offset := int(slot % int64(len(r.buckets)))

	r.mu.Lock()
	defer r.mu.Unlock()
	bucket := &r.buckets[offset]
	if bucket.slot != slot {
		bucket.slot = slot
		bucket.sum = 0
		bucket.count = 0
	}
	bucket.sum += value
	bucket.count++
}

func (r *rollingCounter) reduce(value func(counterBucket) float64, aggregate func(float64, float64) float64, fallback float64) float64 {
	slot := r.currentSlot()
	size := int64(len(r.buckets))
	result := fallback

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, bucket := range r.buckets {
		if bucket.count == 0 || bucket.slot == slot || slot-bucket.slot >= size || bucket.slot > slot {
			continue
		}
		result = aggregate(result, value(bucket))
	}
	return result
}

func (r *rollingCounter) currentSlot() int64 {
	return time.Now().UnixNano() / int64(r.bucketDuration)
}
