package client

import (
	"context"

	"go-common/app/interface/main/playlist/model"
	"go-common/library/net/rpc"
)

const (
	_setStat = "RPC.SetStat"

	_appid = "community.service.playlist"
)

var (
	_noReply = &struct{}{}
)

// Service struct info.
type Service struct {
	client *rpc.Client2
}

// New new servcie instance and return.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// SetStat updates playlist stat cache.
func (s *Service) SetStat(c context.Context, arg *model.ArgStats) (err error) {
	err = s.client.Call(c, _setStat, arg, _noReply)
	return
}
