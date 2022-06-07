package selector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeer(t *testing.T) {
	var p Peer = Peer{
		Node: mockWeightedNode{},
	}
	ctx := NewPeerContext(context.Background(), &p)
	p2, ok := FromPeerContext(ctx)
	if ok {
		assert.NotNil(t, p2.Node)
	}
}

func TestNotPeer(t *testing.T) {
	_, ok := FromPeerContext(context.Background())

	assert.True(t, !ok)
}
