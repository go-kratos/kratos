package v1

import (
	"context"
	"fmt"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// DiscoveryAppID .
const DiscoveryAppID = "main.archive.creative"

// NewClient new grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (CreativeClient, error) {
	client := warden.NewClient(cfg, opts...)
	// cc, err := client.Dial(context.Background(), "127.0.0.1:9000")
	cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", DiscoveryAppID))
	if err != nil {
		return nil, err
	}
	return NewCreativeClient(cc), nil
}
