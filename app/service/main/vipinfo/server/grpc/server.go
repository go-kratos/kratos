// Package server generate by warden_gen
package server

import (
	"context"

	v1 "go-common/app/service/main/vipinfo/api"
	service "go-common/app/service/main/vipinfo/service"
	"go-common/library/net/rpc/warden"
)

// New VipInfo warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	v1.RegisterVipInfoServer(ws.Server(), &server{svr})
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ v1.VipInfoServer = &server{}

// Info get vipinfo by mid.
func (s *server) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoReply, err error) {
	var info *v1.ModelInfo
	if info, err = s.svr.Info(ctx, req.Mid); err != nil {
		return
	}
	return &v1.InfoReply{Res: info}, nil
}

// Infos get vipinfos by mids
func (s *server) Infos(ctx context.Context, req *v1.InfosReq) (res *v1.InfosReply, err error) {
	var infos map[int64]*v1.ModelInfo
	if infos, err = s.svr.Infos(ctx, req.Mids); err != nil {
		return
	}
	return &v1.InfosReply{Res: infos}, nil
}
