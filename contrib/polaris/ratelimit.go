package polaris

import (
	"context"
	"strings"

	"github.com/go-kratos/aegis/ratelimit"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/polarismesh/polaris-go/pkg/model"
)

// ErrLimitExceed is service unavailable due to rate limit exceeded.
var (
	ErrLimitExceed = errors.New(429, "RATELIMIT", "service unavailable due to rate limit exceeded")
)

// Ratelimit Request rate limit middleware
func Ratelimit(l Limiter) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				var args []model.Argument
				headers := tr.RequestHeader()
				// handle header
				for _, header := range headers.Keys() {
					args = append(args, model.BuildHeaderArgument(header, headers.Get(header)))
				}
				// handle http
				if ht, ok := tr.(*http.Transport); ok {
					// url query
					for key, values := range ht.Request().URL.Query() {
						args = append(args, model.BuildQueryArgument(key, strings.Join(values, ",")))
					}
				}
				done, e := l.Allow(tr.Operation(), args...)
				if e != nil {
					// rejected
					return nil, ErrLimitExceed
				}
				// allowed
				reply, err = handler(ctx, req)
				done(ratelimit.DoneInfo{Err: err})
				return
			}
			return reply, nil
		}
	}
}
