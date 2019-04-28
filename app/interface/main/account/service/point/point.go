package point

import (
	"context"

	"go-common/app/interface/main/account/conf"
	"go-common/app/service/main/point/model"
	pointrpc "go-common/app/service/main/point/rpc/client"
)

// Service struct of service.
type Service struct {
	// conf
	c *conf.Config
	// rpc
	pointRPC *pointrpc.Service
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		pointRPC: pointrpc.New(c.RPCClient2.Point),
	}
	return
}

// PointInfo point info.
func (s *Service) PointInfo(c context.Context, mid int64) (res *model.PointInfo, err error) {
	res, err = s.pointRPC.PointInfo(c, &model.ArgRPCMid{Mid: mid})
	return
}

// PointPage point page.
func (s *Service) PointPage(c context.Context, a *model.ArgRPCPointHistory) (res *model.PointHistoryResp, err error) {
	res, err = s.pointRPC.PointHistory(c, a)
	return
}
