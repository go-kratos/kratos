package v1

import (
	"context"
	"go-common/library/net/rpc/warden"
	"google.golang.org/grpc"
)

const kAppID = "live.rtc"

type Client struct {
	RtcClient
}

// NewClient new resource grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (*Client, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+kAppID)
	if err != nil {
		return nil, err
	}
	cli := &Client{}
	cli.RtcClient = NewRtcClient(conn)
	return cli, nil
}
