package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID unique app id for service discovery
const AppID = "archive.service.up"

// NewClient new member grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (UpClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	return NewUpClient(conn), nil
}
