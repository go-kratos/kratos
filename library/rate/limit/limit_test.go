package limit

import (
	"context"
	"testing"
	"time"

	"go-common/library/rate"
)

func worker(qps int64, ch chan struct{}) {
	for {
		<-ch
		time.Sleep(time.Duration(int64(time.Second) / qps))
	}
}

func TestRateSuccess(t *testing.T) {
	ch := make(chan struct{})
	go worker(100, ch)
	failed := producer(New(nil), 100, ch)
	if failed > 0 {
		t.Fatalf("Should be rejected 0 time,but (%d)", failed)
	}
}

func TestRateFail(t *testing.T) {
	ch := make(chan struct{})
	go worker(100, ch)
	failed := producer(New(nil), 200, ch)
	if failed < 900 {
		t.Fatalf("Should be rejected more than 900 times,but (%d)", failed)
	}
}

func TestRateFailMuch(t *testing.T) {
	ch := make(chan struct{})
	go worker(10, ch)
	failed := producer(New(nil), 200, ch)
	if failed < 1600 {
		t.Fatalf("Should be rejected more than 1600 times,but (%d)", failed)
	}
}

func producer(l *Limiter, qps int64, ch chan struct{}) (failed int) {
	for i := 0; i < int(qps)*10; i++ {
		go func() {
			done, err := l.Allow(context.Background())
			defer done(rate.Success)
			if err == nil {
				ch <- struct{}{}
			} else {
				failed++
			}
		}()
		time.Sleep(time.Duration(int64(time.Second) / qps))
	}
	return
}
