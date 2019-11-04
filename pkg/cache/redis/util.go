package redis

import (
	"context"
	"time"
)

func shrinkDeadline(ctx context.Context, timeout time.Duration) time.Time {
	var timeoutTime = time.Now().Add(timeout)
	if ctx == nil {
		return timeoutTime
	}
	if deadline, ok := ctx.Deadline(); ok && timeoutTime.After(deadline) {
		return deadline
	}
	return timeoutTime
}
