package v1

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

// AppID unique app id for service discovery
const AppID = "passport.service.auth"

// NewClient new auth gRPC client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (AuthClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	return NewAuthClient(conn), nil
}
