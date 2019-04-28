package netutil

import (
	"math/rand"
	"time"
)

// DefaultBackoffConfig uses values specified for backoff in common.
var DefaultBackoffConfig = BackoffConfig{
	MaxDelay:  120 * time.Second,
	BaseDelay: 1.0 * time.Second,
	Factor:    1.6,
	Jitter:    0.2,
}

// Backoff defines the methodology for backing off after a call failure.
type Backoff interface {
	// Backoff returns the amount of time to wait before the next retry given
	// the number of consecutive failures.
	Backoff(retries int) time.Duration
}

// BackoffConfig defines the parameters for the default backoff strategy.
type BackoffConfig struct {
	// MaxDelay is the upper bound of backoff delay.
	MaxDelay time.Duration

	// baseDelay is the amount of time to wait before retrying after the first
	// failure.
	BaseDelay time.Duration

	// factor is applied to the backoff after each retry.
	Factor float64

	// jitter provides a range to randomize backoff delays.
	Jitter float64
}

/*
// NOTE TODO avoid use unexcept config.
func (bc *BackoffConfig) Fix() {
	md := bc.MaxDelay
	*bc = DefaultBackoffConfig

	if md > 0 {
		bc.MaxDelay = md
	}
}
*/

// Backoff returns the amount of time to wait before the next retry given
// the number of consecutive failures.
func (bc *BackoffConfig) Backoff(retries int) time.Duration {
	if retries == 0 {
		return bc.BaseDelay
	}
	backoff, max := float64(bc.BaseDelay), float64(bc.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= bc.Factor
		retries--
	}
	if backoff > max {
		backoff = max
	}
	// Randomize backoff delays so that if a cluster of requests start at
	// the same time, they won't operate in lockstep.
	backoff *= 1 + bc.Jitter*(rand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}
