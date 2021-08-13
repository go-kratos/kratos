package selector

import (
	"context"
	"regexp"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type (
	transporter func(ctx context.Context) (transport.Transporter, bool)
)

var (
	// serverTransporter is get server transport.Transporter from ctx
	serverTransporter transporter = func(ctx context.Context) (transport.Transporter, bool) {
		return transport.FromServerContext(ctx)
	}
	// clientTransporter is get client transport.Transporter from ctx
	clientTransporter transporter = func(ctx context.Context) (transport.Transporter, bool) {
		return transport.FromClientContext(ctx)
	}
)

// Builder is a selector builder
type Builder struct {
	client bool

	rootOpMatcherBuilder orMatcherBuilder

	ms []middleware.Middleware
}

// Server selector middleware
func Server(ms ...middleware.Middleware) *Builder {
	return &Builder{ms: ms}
}

// Client selector middleware
func Client(ms ...middleware.Middleware) *Builder {
	return &Builder{client: true, ms: ms}
}

// Prefix is with Builder's prefix
func (b *Builder) Prefix(prefix ...string) *Builder {
	for _, s := range prefix {
		b.rootOpMatcherBuilder.push(prefixMather(s))
	}
	return b
}

// Regex is with Builder's regex
// panics if any expression cannot be parsed.
func (b *Builder) Regex(regex ...string) *Builder {
	for _, s := range regex {
		re := regexp.MustCompile(s)
		b.rootOpMatcherBuilder.push(regexMatcher{re: re})
	}
	return b
}

// Path is with Builder's path
func (b *Builder) Path(path ...string) *Builder {
	for _, s := range path {
		b.rootOpMatcherBuilder.push(pathMather(s))
	}
	return b
}

func (b *Builder) OperationMatcher(matcher ...OperationMatcher) *Builder {
	b.rootOpMatcherBuilder.push(matcher...)
	return b
}

// Build is Builder's Build, for example: Server().Path(m1,m2).Build()
func (b *Builder) Build() middleware.Middleware {
	var transporter func(ctx context.Context) (transport.Transporter, bool)
	if b.client {
		transporter = clientTransporter
	} else {
		transporter = serverTransporter
	}
	return selector(transporter, b.rootOpMatcherBuilder.build(), b.ms...)
}

// selector middleware
func selector(transporter transporter, opMatcher OperationMatcher, ms ...middleware.Middleware) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			info, ok := transporter(ctx)
			if !ok {
				return handler(ctx, req)
			}

			if !opMatcher.Match(info.Operation()) {
				return handler(ctx, req)
			}
			return middleware.Chain(ms...)(handler)(ctx, req)
		}
	}
}
