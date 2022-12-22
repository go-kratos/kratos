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
	ErrLimitExceed         = errors.New(429, "RATELIMIT", "service unavailable due to rate limit exceeded")
	querySplit             = "&"
	paramSplit             = "="
	polarisArgumentsSplist = ","
)

// Server ratelimiter middleware
func Server(l Limiter) middleware.Middleware {
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
					// handle query
					querys := strings.Split(ht.Request().URL.RawQuery, querySplit)
					requestStringMap := make(map[string]string)
					for _, query := range querys {
						// the build result of test=1&test=2 is '1,2'
						params := strings.Split(query, paramSplit)
						if requestStringMap[params[0]] == "" {
							requestStringMap[params[0]] = requestStringMap[params[0]] + params[1]
							continue
						}
						requestStringMap[params[0]] = requestStringMap[params[0]] + polarisArgumentsSplist + params[1]
					}
					for k, v := range requestStringMap {
						args = append(args, model.BuildQueryArgument(k, v))
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
			return nil, errors.New(400, "Error with transport.FromServerContext", "Error with transport.FromServerContext")
		}
	}
}
