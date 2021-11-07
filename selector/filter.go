package selector

import "context"

// Filter is select filter.
type Filter func(context.Context, []Node) []Node

// NodeFilter is node filter.
// If it returns false, the node will be removed out from the balancer pick list
type NodeFilter func(node Node) bool
