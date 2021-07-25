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

type Builder struct {
	client bool

	prefix []string
	regex  []string
	path   []string

	notPrefix []string
	notRegex  []string
	notPath   []string

	ms []middleware.Middleware
}

func Server(ms ...middleware.Middleware) *Builder {
	return &Builder{ms: ms}
}

func Client(ms ...middleware.Middleware) *Builder {
	return &Builder{client: true, ms: ms}
}

func (b *Builder) Prefix(prefix ...string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Regex(regex ...string) *Builder {
	b.regex = regex
	return b
}
func (b *Builder) Path(path ...string) *Builder {
	b.path = path
	return b
}

func (b *Builder) NotPrefix(prefix ...string) *Builder {
	b.notPrefix = prefix
	return b
}

func (b *Builder) NotRegex(regex ...string) *Builder {
	b.notRegex = regex
	return b
}
func (b *Builder) NotPath(path ...string) *Builder {
	b.notPath = path
	return b
}

func (b *Builder) Build() middleware.Middleware {
	var transporter func(ctx context.Context) (transport.Transporter, bool)
	if b.client == true {
		transporter = clientTransporter
	} else {
		transporter = serverTransporter
	}
	return selector(transporter, b.match, b.ms...)
}

func (b *Builder) match(operation string) bool {
	if len(b.notPrefix)+len(b.notRegex)+len(b.notPath) > 0 {
		for _, prefix := range b.notPrefix {
			if prefixMatch(prefix, operation) {
				return false
			}
		}
		for _, regex := range b.notRegex {
			if regexMatch(regex, operation) {
				return false
			}
		}
		for _, path := range b.notPath {
			if pathMatch(path, operation) {
				return false
			}
		}
		return true
	}

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
	return false
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

func pathMatch(path string, operation string) bool {
	return path == operation
}

func prefixMatch(prefix string, operation string) bool {
	return strings.HasPrefix(operation, prefix)
}

func regexMatch(regex string, operation string) bool {
	return regexMatchByString(regex, operation)
}

func regexMatchByString(pattern string, operation string) bool {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return r.FindString(operation) == operation
}
