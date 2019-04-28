package grpc

import (
	"context"
	rpc "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/service"
	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// New 生成grpc服务
func New(srv *service.Service) *warden.Server {
	// conf := &warden.ServerConfig{Addr: "0.0.0.0:9009"}
	// s := warden.NewServer(conf)
	s := warden.NewServer(nil)
	rpc.RegisterRecsysServer(s.Server(), srv)
	s.Use(middleware())
	_, err := s.Start()
	if err != nil {
		panic("run server failed!" + err.Error())
	}
	return s
}

// middleware middleware
func middleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//call chain
		resp, err = handler(ctx, req)
		return
	}
}
