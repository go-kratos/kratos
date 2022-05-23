package selector

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
)

// ErrNoAvailable is no available node.
var ErrNoAvailable = errors.ServiceUnavailable("no_available_node", "")

// Selector is node pick balancer.
type Selector interface {
	Rebalancer

	// Select nodes
	// if err == nil, selected and done must not be empty.
	Select(ctx context.Context, opts ...SelectOption) (selected Node, done DoneFunc, err error)
}

// Rebalancer is nodes rebalancer.
type Rebalancer interface {
	// Apply is apply all nodes when any changes happen
	Apply(nodes []Node)
}

// Builder build selector
type Builder interface {
	Build() Selector
}

// Node is node interface.
type Node interface {
	// Scheme is service node scheme
	Scheme() string

	// Address is the unique address under the same service
	Address() string

	// ServiceName is service name
	ServiceName() string

	// InitialWeight is the initial value of scheduling weight
	// if not set return nil
	InitialWeight() *int64

	// Version is service node version
	Version() string

	// Metadata is the kv pair metadata associated with the service instance.
	// version,namespace,region,protocol etc..
	Metadata() map[string]string
}

// DoneInfo is callback info when RPC invoke done.
type DoneInfo struct {
	// Response Error
	Err error
	// Response Metadata
	ReplyMeta ReplyMeta

	// BytesSent indicates if any bytes have been sent to the server.
	BytesSent bool
	// BytesReceived indicates if any byte has been received from the server.
	BytesReceived bool
}

// ReplyMeta is Reply Metadata.
type ReplyMeta interface {
	Get(key string) string
}

// DoneFunc is callback function when RPC invoke done.
type DoneFunc func(ctx context.Context, di DoneInfo)
