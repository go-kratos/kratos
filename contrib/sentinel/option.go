package sentinel

import "context"

type (
	Option  func(*options)
	options struct {
		resourceExtract func(ctx context.Context, req interface{}) string
		blockFallback   func(ctx context.Context, req interface{}) (interface{}, error)
	}
)

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	for _, opt := range opts {
		opt(optCopy)
	}

	return optCopy
}

// WithResourceExtractor sets the resource extractor of the web requests.
func WithResourceExtractor(fn func(ctx context.Context, req interface{}) string) Option {
	return func(opts *options) {
		opts.resourceExtract = fn
	}
}

// WithBlockFallback sets the fallback handler when requests are blocked.
func WithBlockFallback(fn func(ctx context.Context, req interface{}) (interface{}, error)) Option {
	return func(opts *options) {
		opts.blockFallback = fn
	}
}
