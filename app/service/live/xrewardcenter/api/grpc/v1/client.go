package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID unique app id for service discovery
const AppID = "live.xrewardcenter"

// NewClient new member grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (AnchorRewardClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	//conn, err := client.Dial(context.Background(), "127.0.0.1:9000")

	if err != nil {
		return nil, err
	}
	return NewAnchorRewardClient(conn), nil
}
