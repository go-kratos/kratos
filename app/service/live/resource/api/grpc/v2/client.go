package v2

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

// AppID 应用程序标识
const AppID = "live.resource"

// Client 对外服务接口
type Client struct {
	UserResourceClient
}

// NewClient 用户资源 grpc Client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (*Client, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	cli := &Client{}
	cli.UserResourceClient = NewUserResourceClient(conn)
	return cli, nil
}
