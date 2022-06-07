package selector

import (
	"context"
)

type peerKey struct{}

// Peer contains the information of the peer for an RPC, such as the address
// and authentication information.
type Peer struct {
	// node is the peer node.
	Node Node
}

// NewPeerContext creates a new context with peer information attached.
func NewPeerContext(ctx context.Context, p *Peer) context.Context {
	return context.WithValue(ctx, peerKey{}, p)
}

// FromPeerContext returns the peer information in ctx if it exists.
func FromPeerContext(ctx context.Context) (p *Peer, ok bool) {
	p, ok = ctx.Value(peerKey{}).(*Peer)
	return
}
