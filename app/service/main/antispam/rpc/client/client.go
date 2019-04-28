package rpc

import (
	"context"

	"go-common/app/service/main/antispam/model"
	"go-common/library/net/rpc"
)

const (
	_appid = "antispam.service"
)

// Client .
type Client struct {
	*rpc.Client2
}

// NewClient .
func NewClient(c *rpc.ClientConfig) *Client {
	s := &Client{}
	s.Client2 = rpc.NewDiscoveryCli(_appid, c)
	return s
}

// Filter .
func (cli *Client) Filter(ctx context.Context, arg *model.Suspicious) (res *model.SuspiciousResp, err error) {
	err = cli.Call(ctx, "Filter.Check", arg, &res)
	return
}
