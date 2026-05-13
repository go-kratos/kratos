package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

func TestBreakerAllowsBeforeRequestThreshold(t *testing.T) {
	breaker := NewBreaker(WithRequest(3), WithWindow(time.Second), WithBucket(2))
	breaker.MarkFailed()
	breaker.MarkFailed()

	if err := breaker.Allow(); err != nil {
		t.Fatalf("Allow() error = %v, want nil", err)
	}
}

func TestBreakerRejectsAfterFailures(t *testing.T) {
	breaker := NewBreaker(WithRequest(1), WithWindow(time.Second), WithBucket(2)).(*Breaker)
	breaker.random = func() float64 { return 0 }
	breaker.MarkFailed()

	if err := breaker.Allow(); !errors.Is(err, ErrNotAllowed) {
		t.Fatalf("Allow() error = %v, want %v", err, ErrNotAllowed)
	}
}

func TestBreakerDefaultRequestThreshold(t *testing.T) {
	breaker := NewBreaker(WithWindow(time.Second), WithBucket(2)).(*Breaker)
	breaker.random = func() float64 { return 0 }

	for range 19 {
		breaker.MarkFailed()
	}
	if err := breaker.Allow(); err != nil {
		t.Fatalf("Allow() error = %v, want nil before default request threshold", err)
	}

	breaker.MarkFailed()
	if err := breaker.Allow(); !errors.Is(err, ErrNotAllowed) {
		t.Fatalf("Allow() error = %v, want %v at default request threshold", err, ErrNotAllowed)
	}
}

func TestBreakerDefaultFailureRatio(t *testing.T) {
	breaker := NewBreaker(WithRequest(20), WithWindow(time.Second), WithBucket(2)).(*Breaker)
	breaker.random = func() float64 { return 0 }

	for range 10 {
		breaker.MarkSuccess()
		breaker.MarkFailed()
	}
	if err := breaker.Allow(); err != nil {
		t.Fatalf("Allow() error = %v, want nil at default failure ratio threshold", err)
	}

	breaker.MarkFailed()
	if err := breaker.Allow(); !errors.Is(err, ErrNotAllowed) {
		t.Fatalf("Allow() error = %v, want %v above default failure ratio threshold", err, ErrNotAllowed)
	}
}

func TestBreakerClosesAfterSuccesses(t *testing.T) {
	breaker := NewBreaker(WithRequest(1), WithFailureRatio(0.5), WithWindow(time.Second), WithBucket(2))
	breaker.MarkSuccess()
	breaker.MarkSuccess()

	if err := breaker.Allow(); err != nil {
		t.Fatalf("Allow() error = %v, want nil", err)
	}
}
