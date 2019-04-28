package v1

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

// AppID unique app id for service discovery
const AppID = "account.service.spy"

// NewClient new identify grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (SpyClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	return NewSpyClient(conn), nil
}
