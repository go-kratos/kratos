package grpc

import (
	pb "go-common/app/service/live/resource/api/grpc/v1"
	v2pb "go-common/app/service/live/resource/api/grpc/v2"
	"go-common/app/service/live/resource/conf"
	svc "go-common/app/service/live/resource/service/v1"
	v2svc "go-common/app/service/live/resource/service/v2"
	"go-common/library/net/rpc/warden"
)

// New
func New(c *conf.Config) *warden.Server {
	ws := warden.NewServer(nil)
	pb.RegisterResourceServer(ws.Server(), svc.NewResourceService(c))
	pb.RegisterSplashServer(ws.Server(), svc.NewSplashService(c))
	pb.RegisterBannerServer(ws.Server(), svc.NewBannerService(c))
	pb.RegisterLiveCheckServer(ws.Server(), svc.NewLiveCheckService(c))
	pb.RegisterTitansServer(ws.Server(), svc.NewTitansService(c))
	v2pb.RegisterUserResourceServer(ws.Server(), v2svc.NewUserResourceService(c))
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
