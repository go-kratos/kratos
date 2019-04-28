package v0

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

const AppID = "live.daoanchor"

type Client struct {
	CreateDataClient
}

// NewClient new anchor grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (*Client, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	cli := &Client{}
	cli.CreateDataClient = NewCreateDataClient(conn)
	return cli, nil
}
