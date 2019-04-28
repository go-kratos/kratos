package grpc

import (
	uprpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/service"
	"go-common/library/net/rpc/warden"
)

// New new a grpc server.
func New(cfg *warden.ServerConfig, s *service.Service) *warden.Server {
	grpc := warden.NewServer(cfg)
	uprpc.RegisterUpServer(grpc.Server(), s)
	grpc, err := grpc.Start()
	if err != nil {
		panic(err)
	}
	return grpc
}
