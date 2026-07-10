package ratelimit

import (
	"errors"
	"testing"
	"time"
)

func TestBBRAllowRecordsDone(t *testing.T) {
	limiter := NewLimiter(WithWindow(time.Second), WithBucket(10))
	done, err := limiter.Allow()
	if err != nil {
		t.Fatalf("Allow() error = %v, want nil", err)
	}
	done(DoneInfo{})

	stat := limiter.Stat()
	if stat.InFlight != 0 {
		t.Fatalf("InFlight = %d, want 0", stat.InFlight)
	}
	if stat.MaxPass == 0 {
		t.Fatalf("MaxPass = %d, want > 0", stat.MaxPass)
	}
}

func TestBBRDropsWhenCPUThresholdExceeded(t *testing.T) {
	limiter := NewLimiter(WithWindow(time.Second), WithBucket(10), WithCPUThreshold(-1))
	limiter.cpu = func() int64 { return 1000 }

	done1, err := limiter.Allow()
	if err != nil {
		t.Fatalf("first Allow() error = %v, want nil", err)
	}
	defer done1(DoneInfo{})
	done2, err := limiter.Allow()
	if err != nil {
		t.Fatalf("second Allow() error = %v, want nil", err)
	}
	defer done2(DoneInfo{})

	if _, err = limiter.Allow(); !errors.Is(err, ErrLimitExceed) {
		t.Fatalf("third Allow() error = %v, want %v", err, ErrLimitExceed)
	}
}
