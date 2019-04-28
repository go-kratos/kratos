package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// .
const (
	DiscoveryID = "push.service.broadcast"
)

// NewClient .
func NewClient(conf *warden.ClientConfig, opts ...grpc.DialOption) (ZergClient, error) {
	client := warden.NewClient(conf, opts...)
	cc, err := client.Dial(context.Background(), "discovery://default/"+DiscoveryID)
	if err != nil {
		return nil, err
	}
	return NewZergClient(cc), nil
}
