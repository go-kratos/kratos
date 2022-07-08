package ratelimit

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kratos/aegis/ratelimit"
)

type (
	ratelimitMock struct {
		reached bool
	}
	ratelimitReachedMock struct {
		reached bool
	}
)

func (r *ratelimitMock) Allow() (ratelimit.DoneFunc, error) {
	return func(_ ratelimit.DoneInfo) {
		r.reached = true
	}, nil
}

func (r *ratelimitReachedMock) Allow() (ratelimit.DoneFunc, error) {
	return func(_ ratelimit.DoneInfo) {
		r.reached = true
	}, errors.New("errored")
}

func Test_WithLimiter(t *testing.T) {
	o := options{
		limiter: &ratelimitMock{},
	}

	WithLimiter(nil)(&o)
	if o.limiter != nil {
		t.Error("The limiter property must be updated.")
	}
}

func Test_Server(t *testing.T) {
	nextValid := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "Hello valid", nil
	}

	rlm := &ratelimitMock{}
	rlrm := &ratelimitReachedMock{}

	_, _ = Server(func(o *options) {
		o.limiter = rlm
	})(nextValid)(context.Background(), nil)
	if !rlm.reached {
		t.Error("The ratelimit must run the done function.")
	}

	_, _ = Server(func(o *options) {
		o.limiter = rlrm
	})(nextValid)(context.Background(), nil)
	if rlrm.reached {
		t.Error("The ratelimit must not run the done function and should be denied.")
	}
}
