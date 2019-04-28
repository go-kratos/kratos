package v1

import (
	"context"
	"go-common/library/net/rpc/warden"
	"google.golang.org/grpc"
)

// AppID 服务app_id
const AppID = "live.xcaptcha"

// Client grpc xcaptcha
type Client struct {
	XCaptchaClient
}

// NewClient new resource grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (*Client, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	cli := &Client{}
	cli.XCaptchaClient = NewXCaptchaClient(conn)
	return cli, nil
}
