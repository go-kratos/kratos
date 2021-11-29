package selector

import "context"

// Filter is select filter.
type Filter func(context.Context, []Node) []Node
