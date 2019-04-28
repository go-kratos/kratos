package client

import (
	"context"

	"go-common/app/service/main/seq-server/model"
	"go-common/library/net/rpc"
)

const (
	_ID      = "RPC.ID"
	_IDInt32 = "RPC.ID32"
)

const (
	_appid = "seq.server"
)

// Service2 is seq rpc client.
type Service2 struct {
	client *rpc.Client2
}

// New2 new a seq rpc client.
func New2(c *rpc.ClientConfig) (s *Service2) {
	s = &Service2{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// ID get id.
func (s *Service2) ID(c context.Context, arg *model.ArgBusiness) (id int64, err error) {
	err = s.client.Call(c, _ID, arg, &id)
	return
}

// ID32 get id32.
func (s *Service2) ID32(c context.Context, arg *model.ArgBusiness) (id int32, err error) {
	err = s.client.Call(c, _IDInt32, arg, &id)
	return
}
