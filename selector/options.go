package selector

// SelectOptions is Select Options.
type SelectOptions struct {
	NodeFilters []NodeFilter
}

// SelectOption is Selector option.
type SelectOption func(*SelectOptions)

// WithNodeFilter with filter options
func WithNodeFilter(fn ...NodeFilter) SelectOption {
	return func(opts *SelectOptions) {
		opts.NodeFilters = fn
	}
}
