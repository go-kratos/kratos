package selector

import (
	"context"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type (
	transporter func(ctx context.Context) (transport.Transporter, bool)
	match       func(operation string) bool
	MatchFunc   match
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

	prefix    []string
	regex     []string
	path      []string
	matchFunc MatchFunc

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
	b.prefix = prefix
	return b
}

// Regex is with Builder's regex
func (b *Builder) Regex(regex ...string) *Builder {
	b.regex = regex
	return b
}

// Path is with Builder's path
func (b *Builder) Path(path ...string) *Builder {
	b.path = path
	return b
}

// Match is with Builder's matchFunc
func (b *Builder) Match(fn MatchFunc) *Builder {
	b.matchFunc = fn
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
	return selector(transporter, b.match, b.ms...)
}

// match is match operation compliance Builder
func (b *Builder) match(operation string) bool {
	for _, prefix := range b.prefix {
		if prefixMatch(prefix, operation) {
			return true
		}
	}
	for _, regex := range b.regex {
		if regexMatch(regex, operation) {
			return true
		}
	}
	for _, path := range b.path {
		if pathMatch(path, operation) {
			return true
		}
	}

	if b.matchFunc != nil {
		return b.matchFunc(operation)
	}
	return false
}

// selector middleware
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

func pathMatch(path string, operation string) bool {
	return path == operation
}

func prefixMatch(prefix string, operation string) bool {
	return strings.HasPrefix(operation, prefix)
}

func regexMatch(regex string, operation string) bool {
	r, err := regexp.Compile(regex)
	if err != nil {
		return false
	}
	return r.FindString(operation) == operation
}
