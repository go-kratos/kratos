package selector

// SelectOptions is Select Options.
type SelectOptions struct {
	Filters []Filter
}

// SelectOption is Selector option.
type SelectOption func(*SelectOptions)

// WithFilter with filter options
func WithFilter(fn ...Filter) SelectOption {
	return func(opts *SelectOptions) {
		opts.Filters = fn
	}
}
