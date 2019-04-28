package client

import (
	"context"

	"go-common/app/service/main/passport/model"
	"go-common/library/net/rpc"
)

const (
	_appid = "passport.service"
)

// Client2 struct
type Client2 struct {
	client *rpc.Client2
}

// New Client2 init
func New(c *rpc.ClientConfig) (s *Client2) {
	s = &Client2{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

const (
	_loginLogs = "RPC.LoginLogs"
)

// LoginLogs get the latest limit login logs.
func (c2 *Client2) LoginLogs(c context.Context, arg *model.ArgLoginLogs) (res []*model.LoginLog, err error) {
	res = make([]*model.LoginLog, 0)
	err = c2.client.Call(c, _loginLogs, arg, &res)
	return
}
