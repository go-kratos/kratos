package balancer

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
)

// ErrNoAvailable is no available node
var ErrNoAvailable = errors.ServiceUnavailable("no_available_node", "")

// Selector is node pick balancer
type Selector interface {
	// Select nodes
	// if err == nil, selected and done must not be empty.
	Select(ctx context.Context, nodes []Node) (selected Node, done Done, err error)
}
