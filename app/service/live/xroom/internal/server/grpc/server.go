package grpc

import (
	pb "go-common/app/service/live/xroom/api"
	"go-common/app/service/live/xroom/internal/service"
	"go-common/library/conf/paladin"
	"go-common/library/net/rpc/warden"
)

// New new a grpc server.
func New(svc *service.Service) *warden.Server {
	var rc struct {
		Server *warden.ServerConfig
	}
	if err := paladin.Get("grpc.toml").UnmarshalTOML(&rc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	ws := warden.NewServer(rc.Server)
	pb.RegisterRoomServer(ws.Server(), svc)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
