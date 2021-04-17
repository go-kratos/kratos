package ratelimit

import (
	"context"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestLimit(t *testing.T) {
	var (
		// rate is 1, burst size is 1
		l1 = New(1, 1)
		// rate is 2, burst size is 2
		l2 = New(2, 2)
		// rate is 1, burst size is 1
		l3 = New(1, 1)

		h = func(_ context.Context, _ interface{}) (interface{}, error) { return nil, nil }
	)
	tests := []struct {
		limter  *rate.Limiter
		sleep   time.Duration
		wantErr bool
	}{
		{l1, 0, false},
		{l1, 0, true}, // not enought tokens
		{l2, 0, false},
		{l2, 0, false},
		{l3, 0, false},
		{l3, time.Second, false}, // after sleep 1s, tokens are enough, so no want error
	}
	for i, tt := range tests {
		if tt.sleep != 0 {
			time.Sleep(tt.sleep)
		}
		if _, err := Limit(tt.limter)(h)(context.Background(), nil); tt.wantErr != (err != nil) {
			t.Errorf("step(%d) Limit() = %v, wantErr = %v", i, err, tt.wantErr)
		}
	}
}
