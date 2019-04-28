package v1

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

const AppID = "live.resource"

type Client struct {
	ResourceClient
	SplashClient
	BannerClient
	LiveCheckClient
	TitansClient
}

// NewClient new resource grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (*Client, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	cli := &Client{}
	cli.ResourceClient = NewResourceClient(conn)
	cli.SplashClient = NewSplashClient(conn)
	cli.BannerClient = NewBannerClient(conn)
	cli.LiveCheckClient = NewLiveCheckClient(conn)
	cli.TitansClient = NewTitansClient(conn)
	return cli, nil
}
