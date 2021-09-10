package selector

import "context"

// SelectOptions is Select Options.
type SelectOptions struct {
	Filters []Filter
}

// SelectOption is Selector option.
type SelectOption func(*SelectOptions)

// Filter is node filter function.
type Filter func(context.Context, []Node) []Node

// WithFilters with filter options
func WithFilters(filters ...Filter) SelectOption {
	return func(opts *SelectOptions) {
		opts.Filters = filters
	}
}
