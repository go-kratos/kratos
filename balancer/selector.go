package balancer

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
)

// ErrNoAvaliable is no avaliable node
var ErrNoAvaliable = errors.ServiceUnavailable("no_avaliable_node", "")

// Selector is node pick balancer
type Selector interface {
	// Select nodes
	// if err == nil, selected must not be empty.
	Select(ctx context.Context, nodes []Node) (selected Node, err error)
}
