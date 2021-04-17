package ratelimit

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"golang.org/x/time/rate"
)

func New(r float64, b int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(r), b)
}

func Limit(l *rate.Limiter) middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if !l.Allow() {
				return nil, errors.ResourceExhausted("Rate limit exceeded", "Rate limit exceeded")
			}
			return h(ctx, req)
		}
	}
}
