package ratelimit

import "errors"

// ErrLimitExceed is returned when the rate limiter rejects a request.
var ErrLimitExceed = errors.New("rate limit exceeded")

// DoneFunc records request completion.
type DoneFunc func(DoneInfo)

// DoneInfo contains request completion metadata.
type DoneInfo struct {
	Err error
}

// Limiter is a rate limiter.
type Limiter interface {
	Allow() (DoneFunc, error)
}
