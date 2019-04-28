package api

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

const AppID = "live.xroomfeed"

type Client struct {
	RecPoolClient RecPoolClient
}

// NewClient xroom-feed grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (*Client, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	cli := &Client{}
	cli.RecPoolClient = NewRecPoolClient(conn)
	return cli, nil
}
