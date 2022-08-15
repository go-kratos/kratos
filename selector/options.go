package selector

// SelectOptions is Select Options.
type SelectOptions struct {
	Filters []NodeFilter
}

// SelectOption is Selector option.
type SelectOption func(*SelectOptions)

// WithFilter with filter options
func WithFilter(fn ...NodeFilter) SelectOption {
	return func(opts *SelectOptions) {
		opts.Filters = fn
	}
}
