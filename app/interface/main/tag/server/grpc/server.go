package grpc

import (
	pb "go-common/app/interface/main/tag/api"
	"go-common/app/interface/main/tag/service"
	"go-common/library/net/rpc/warden"
)

type grpcServer struct {
	svr *service.Service
}

// New new a grpc server.
func New(cfg *warden.ServerConfig, svr *service.Service) *warden.Server {
	grpc := warden.NewServer(cfg)
	pb.RegisterTagRPCServer(grpc.Server(), &grpcServer{svr: svr})
	grpc, err := grpc.Start()
	if err != nil {
		panic(err)
	}
	return grpc
}
