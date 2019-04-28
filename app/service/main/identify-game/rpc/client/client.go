package client

import (
	"context"

	"go-common/app/service/main/identify-game/model"
	"go-common/library/net/rpc"
)

const (
	_delCache = "RPC.DelCache"
	_appid    = "identify.service.game"
)

var (
	_noRes = &struct{}{}
)

// Client Request Client
type Client struct {
	client *rpc.Client2
}

// New Request Client
func New(c *rpc.ClientConfig) (cli *Client) {
	cli = &Client{
		client: rpc.NewDiscoveryCli(_appid, c),
	}
	return
}

// DelCache del token cache.
func (cli *Client) DelCache(c context.Context, arg *model.CleanCacheArgs) (err error) {
	err = cli.client.Call(c, _delCache, arg, _noRes)
	return
}
