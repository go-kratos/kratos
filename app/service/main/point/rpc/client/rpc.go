package client

import (
	"context"

	"go-common/app/service/main/point/model"
	"go-common/library/net/rpc"
)

const (
	_pointInfo    = "RPC.PointInfo"
	_pointConsume = "RPC.ConsumePoint"
	_pointAdd     = "RPC.PointAdd"
	_pointHistory = "RPC.PointHistory"
	_pointAddByBp = "RPC.PointAddByBp"
)

const (
	_appid = "account.service.point"
)

// Service is a question service.
type Service struct {
	client *rpc.Client2
}

// New new rpc service.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Ping rpc ping.
func (s *Service) Ping(c context.Context, arg *struct{}) (res *int, err error) {
	err = s.client.Call(c, "RPC.Ping", arg, res)
	return
}

//PointInfo point info.
func (s *Service) PointInfo(c context.Context, arg *model.ArgRPCMid) (res *model.PointInfo, err error) {
	res = new(model.PointInfo)
	err = s.client.Call(c, _pointInfo, arg, res)
	return
}

//ConsumePoint consume point.
func (s *Service) ConsumePoint(c context.Context, arg *model.ArgPointConsume) (status int8, err error) {
	err = s.client.Call(c, _pointConsume, arg, &status)
	return
}

//AddPoint add point.
func (s *Service) AddPoint(c context.Context, arg *model.ArgPoint) (status int8, err error) {
	err = s.client.Call(c, _pointAdd, arg, &status)
	return
}

// PointHistory point history.
func (s *Service) PointHistory(c context.Context, arg *model.ArgRPCPointHistory) (res *model.PointHistoryResp, err error) {
	res = new(model.PointHistoryResp)
	err = s.client.Call(c, _pointHistory, arg, res)
	return
}

// PointAddByBp point add by bp.
func (s *Service) PointAddByBp(c context.Context, arg *model.ArgPointAdd) (res *int64, err error) {
	err = s.client.Call(c, _pointAddByBp, arg, &res)
	return
}
