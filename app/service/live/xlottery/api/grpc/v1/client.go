package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID unique app id for service discovery
const AppID = "live.xlottery"

// Client grpc client for interface
type Client struct {
	CapsuleClient
	StormClient
}

// NewClient new lottery grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (*Client, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	cli := &Client{}
	cli.CapsuleClient = NewCapsuleClient(conn)
	cli.StormClient = NewStormClient(conn)
	return cli, nil
}
