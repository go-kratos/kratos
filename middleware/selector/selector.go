package selector

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"regexp"
	"strings"
)

func MatchPrefix(prefix string, ms ...middleware.Middleware) middleware.Middleware {
	return selector(func(operation string) bool {
		return strings.HasPrefix(operation, prefix)
	}, ms...)
}

func MatchRegex(pattern string, ms ...middleware.Middleware) middleware.Middleware {
	return selector(func(operation string) bool {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return false
		}
		return r.FindString(operation) == operation
	}, ms...)

}

func MatchFull(route string, ms ...middleware.Middleware) middleware.Middleware {
	return selector(func(operation string) bool {
		return route == operation
	}, ms...)
}

func selector(match func(operation string) bool, ms ...middleware.Middleware) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			info, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}
			operation := info.Operation()
			if !match(operation) {
				return handler(ctx, req)
			}
			return middleware.Chain(ms...)(handler)(ctx, req)
		}
	}
}
