package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShrinkDeadline(t *testing.T) {
	t.Run("test not deadline", func(t *testing.T) {
		timeout := time.Second
		timeoutTime := time.Now().Add(timeout)
		tm := shrinkDeadline(context.Background(), timeout)
		assert.True(t, tm.After(timeoutTime))
	})
	t.Run("test big deadline", func(t *testing.T) {
		timeout := time.Second
		timeoutTime := time.Now().Add(timeout)
		deadlineTime := time.Now().Add(2 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		tm := shrinkDeadline(ctx, timeout)
		assert.True(t, tm.After(timeoutTime) && tm.Before(deadlineTime))
	})
	t.Run("test small deadline", func(t *testing.T) {
		timeout := time.Second
		deadlineTime := time.Now().Add(500 * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		tm := shrinkDeadline(ctx, timeout)
		assert.True(t, tm.After(deadlineTime) && tm.Before(time.Now().Add(timeout)))
	})
}
