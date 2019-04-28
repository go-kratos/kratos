package grpc

import (
	pb "go-common/app/service/bbq/topic/api"
	"go-common/app/service/bbq/topic/internal/service"
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
	pb.RegisterTopicServer(ws.Server(), svc)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
