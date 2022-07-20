package retry

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Backoff defines constract for backoff strategies
type Backoff interface {
	Next(retry uint) time.Duration
}

// constantBackoff is a constant backoff
type constantBackoff struct {
	initialInterval       int64
	maximumJitterInterval int64
}

func (cb *constantBackoff) Next(retry uint) time.Duration {
	return (time.Duration(cb.initialInterval) * time.Millisecond) + (time.Duration(rand.Int63n(cb.maximumJitterInterval+1)) * time.Millisecond)
}

func NewConstantBackoff(initialInterval, maximumJitterInterval time.Duration) Backoff {
	if maximumJitterInterval < 0 {
		maximumJitterInterval = 0
	}
	return &constantBackoff{
		initialInterval:       int64(initialInterval / time.Millisecond),
		maximumJitterInterval: int64(maximumJitterInterval / time.Millisecond),
	}
}

// exponentialBackoff is a exponential backoff
type exponentialBackoff struct {
	exponentFactor        float64
	initialInterval       float64
	maxInterval           float64
	maximumJitterInterval int64
}

func (eb *exponentialBackoff) Next(retry uint) time.Duration {
	// https://aws.amazon.com/cn/blogs/architecture/exponential-backoff-and-jitter/
	return time.Duration(math.Min(eb.initialInterval*math.Pow(eb.exponentFactor, float64(retry)), eb.maxInterval)+
		float64(rand.Int63n(eb.maximumJitterInterval+1))) * time.Millisecond
}

func NewExponentialBackoff(exponentFactor float64, initialInterval, maxInterval, maximumJitterInterval time.Duration) Backoff {
	if maximumJitterInterval < 0 {
		maximumJitterInterval = 0
	}
	return &exponentialBackoff{
		exponentFactor:        exponentFactor,
		initialInterval:       float64(initialInterval / time.Millisecond),
		maxInterval:           float64(maxInterval / time.Millisecond),
		maximumJitterInterval: int64(maximumJitterInterval / time.Millisecond),
	}
}

// Retriable implements retriers
type Retriable interface {
	NextInterval(retry uint) time.Duration
}

// RetriableFunc is an adapter to allow the use of ordinary functions
// as a Retriable
type RetriableFunc func(retry uint) time.Duration

func (f RetriableFunc) NextInterval(retry uint) time.Duration {
	return f(retry)
}

func NewRetrierFunc(f RetriableFunc) Retriable {
	return f
}

// retrier with some backoff strategy
type retrier struct {
	backoff Backoff
}

func (r *retrier) NextInterval(retry uint) time.Duration {
	return r.backoff.Next(retry)
}

func NewRetrier(backoff Backoff) Retriable {
	return &retrier{
		backoff: backoff,
	}
}

// a null object for retriable
type noRetrier struct{}

func (r *noRetrier) NextInterval(retry uint) time.Duration {
	return 0 * time.Millisecond
}

func NewNoRetrier() Retriable {
	return &noRetrier{}
}

// Strategy is retry strategy
type Strategy struct {
	Attempts   int
	Retrier    Retriable
	Conditions []Condition
}

func NewStrategy(attempts int, retrier Retriable, conditions []Condition) *Strategy {
	return &Strategy{Attempts: attempts, Retrier: retrier, Conditions: conditions}
}

func (s *Strategy) JudgeConditions(r Resp) bool {
	for _, cond := range s.Conditions {
		if cond.Judge(r) {
			return true
		}
	}
	return false
}
