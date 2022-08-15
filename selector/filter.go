package selector

import "context"

// NodeFilter is select filter.
type NodeFilter func(context.Context, []Node) []Node
