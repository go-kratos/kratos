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
// param prefix's type is string or []string{}
func ServerMatchPrefix(prefix interface{}, ms ...middleware.Middleware) middleware.Middleware {
	return selector(serverTransporter, func(operation string) bool {
		return prefixMatch(prefix, operation)
	}, ms...)
}

// ServerMatchRegex is match specific routes by regex
// param pattern's type is string or []string{}
func ServerMatchRegex(pattern interface{}, ms ...middleware.Middleware) middleware.Middleware {
	return selector(serverTransporter, func(operation string) bool {
		return regexMatch(pattern, operation)
	}, ms...)

}

// ServerMatchFull is match specific routes
// param route's type is string or []string{}
func ServerMatchFull(route interface{}, ms ...middleware.Middleware) middleware.Middleware {
	return selector(serverTransporter, func(operation string) bool {
		return fullMatch(route, operation)
	}, ms...)
}

// ClientMatchPrefix is match specific routes by prefix
// param prefix's type is string or []string{}
func ClientMatchPrefix(prefix interface{}, ms ...middleware.Middleware) middleware.Middleware {
	return selector(clientTransporter, func(operation string) bool {
		return prefixMatch(prefix, operation)
	}, ms...)
}

// ClientMatchRegex is match specific routes by regex
// param pattern's type is string or []string{}
func ClientMatchRegex(pattern interface{}, ms ...middleware.Middleware) middleware.Middleware {
	return selector(clientTransporter, func(operation string) bool {
		return regexMatch(pattern, operation)
	}, ms...)

}

// ClientMatchFull is match specific routes
// param route's type is string or []string{}
func ClientMatchFull(route interface{}, ms ...middleware.Middleware) middleware.Middleware {
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

func fullMatch(route interface{}, operation string) bool {
	switch v := route.(type) {
	case string:
		return v == operation
	case []string:
		for _, s := range v {
			if s == operation {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func prefixMatch(prefix interface{}, operation string) bool {
	switch v := prefix.(type) {
	case string:
		return strings.HasPrefix(operation, v)
	case []string:
		for _, s := range v {
			if strings.HasPrefix(operation, s) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func regexMatch(pattern interface{}, operation string) bool {
	switch v := pattern.(type) {
	case string:
		return regexMatchByString(v, operation)
	case []string:
		for _, s := range v {
			if regexMatchByString(s, operation) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func regexMatchByString(pattern string, operation string) bool {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return r.FindString(operation) == operation
}
