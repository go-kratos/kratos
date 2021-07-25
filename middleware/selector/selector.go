package selector

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"regexp"
	"strings"
)

type (
	transporter func(ctx context.Context) (transport.Transporter, bool)
	match       func(operation string) bool
)

var (
	serverTransporter transporter = func(ctx context.Context) (transport.Transporter, bool) {
		return transport.FromServerContext(ctx)
	}
	clientTransporter transporter = func(ctx context.Context) (transport.Transporter, bool) {
		return transport.FromClientContext(ctx)
	}
)

// ServerMatchPrefix is match specific routes by prefix
func ServerMatchPrefix(prefix string, ms ...middleware.Middleware) middleware.Middleware {
	return selector(serverTransporter, func(operation string) bool {
		return prefixMatch(prefix, operation)
	}, ms...)
}

// ServerMatchRegex is match specific routes by regex
func ServerMatchRegex(pattern string, ms ...middleware.Middleware) middleware.Middleware {
	return selector(serverTransporter, func(operation string) bool {
		return regexMatch(pattern, operation)
	}, ms...)

}

// ServerMatchFull is match specific routes
func ServerMatchFull(route string, ms ...middleware.Middleware) middleware.Middleware {
	return selector(serverTransporter, func(operation string) bool {
		return fullMatch(route, operation)
	}, ms...)
}

// ClientMatchPrefix is match specific routes by prefix
func ClientMatchPrefix(prefix string, ms ...middleware.Middleware) middleware.Middleware {
	return selector(clientTransporter, func(operation string) bool {
		return prefixMatch(prefix, operation)
	}, ms...)
}

// ClientMatchRegex is match specific routes by regex
func ClientMatchRegex(pattern string, ms ...middleware.Middleware) middleware.Middleware {
	return selector(clientTransporter, func(operation string) bool {
		return regexMatch(pattern, operation)
	}, ms...)

}

// ClientMatchFull is match specific routes
func ClientMatchFull(route string, ms ...middleware.Middleware) middleware.Middleware {
	return selector(clientTransporter, func(operation string) bool {
		return fullMatch(route, operation)
	}, ms...)
}

func selector(transporter transporter, match match, ms ...middleware.Middleware) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {

			info, ok := transporter(ctx)
			if !ok {
				return handler(ctx, req)
			}

			if !match(info.Operation()) {
				return handler(ctx, req)
			}
			return middleware.Chain(ms...)(handler)(ctx, req)
		}
	}
}

func fullMatch(route string, operation string) bool {
	return route == operation
}

func prefixMatch(prefix string, operation string) bool {
	return strings.HasPrefix(operation, prefix)
}

func regexMatch(pattern string, operation string) bool {
	return regexMatchByString(pattern, operation)
}

func regexMatchByString(pattern string, operation string) bool {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return r.FindString(operation) == operation
}
