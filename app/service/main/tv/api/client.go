package api

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

// AppID unique app id for service discovery
const AppID = "main.web-svr.tv-service"

// NewClient new tv vip grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (TVServiceClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	return NewTVServiceClient(conn), nil
}
