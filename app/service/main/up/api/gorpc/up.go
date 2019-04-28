package client

import (
	"context"

	"go-common/app/service/main/up/model"
	"go-common/library/net/rpc"
)

const (
	_Special     = "RPC.Special"
	_Info        = "RPC.Info"
	_UpStatBase  = "RPC.UpStatBase"
	_SetUpSwitch = "RPC.SetUpSwitch"
	_UpSwitch    = "RPC.UpSwitch"
)
const (
	_appid = "archive.service.up"
)

var (
// _noRes = &struct{}{}
)

// Service rpc client
type Service struct {
	client *rpc.Client2
}

// RPC rpc
type RPC interface {
	// DEPRECATED: Please use gRPC service of func UpGroupMids instead, but must get datas in many times by one time of max 1000.
	Special(c context.Context, arg *model.ArgSpecial) (res []*model.UpSpecial, err error)
	// DEPRECATED: Please use gRPC service of func UpAttr instead.
	Info(c context.Context, arg *model.ArgInfo) (res *model.UpInfo, err error)
}

// New create
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Special special
// DEPRECATED: Please use gRPC service of func UpGroupMids instead, but must get datas in many times by one time of max 1000.
func (s *Service) Special(c context.Context, arg *model.ArgSpecial) (res []*model.UpSpecial, err error) {
	err = s.client.Call(c, _Special, arg, &res)
	return
}

// Info info
// DEPRECATED: Please use gRPC service of func UpAttr instead.
func (s *Service) Info(c context.Context, arg *model.ArgInfo) (res *model.UpInfo, err error) {
	err = s.client.Call(c, _Info, arg, &res)
	return
}

// UpStatBase base statis
// DEPRECATED: Please use gRPC service func UpBaseStats instead.
func (s *Service) UpStatBase(c context.Context, arg *model.ArgMidWithDate) (res *model.UpBaseStat, err error) {
	err = s.client.Call(c, _UpStatBase, arg, &res)
	return
}

// SetUpSwitch set up switch
// DEPRECATED: Please use gRPC service func SetUpSwitch instead.
func (s *Service) SetUpSwitch(c context.Context, arg *model.ArgUpSwitch) (res *model.PBSetUpSwitchRes, err error) {
	err = s.client.Call(c, _SetUpSwitch, arg, &res)
	return
}

// UpSwitch get up switch
// DEPRECATED: Please use gRPC service func UpSwitch instead.
func (s *Service) UpSwitch(c context.Context, arg *model.ArgUpSwitch) (res *model.PBUpSwitch, err error) {
	err = s.client.Call(c, _UpSwitch, arg, &res)
	return
}
