package v1

import (
	"context"
	"fmt"

	"go-common/library/net/rpc/warden"
	"google.golang.org/grpc"
)

// DiscoveryAppID .
const DiscoveryAppID = "video.live.streamng"

// NewClient 建立grpc连接
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (StreamClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", DiscoveryAppID))

	if err != nil {
		return nil, err
	}

	return NewStreamClient(cc), nil
}
