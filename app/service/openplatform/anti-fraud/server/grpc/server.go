package grpc

import (
	"context"
	"go-common/app/service/openplatform/anti-fraud/service"
	"go-common/library/net/rpc/warden"
	"google.golang.org/grpc"
)

//New 生成rpc服务
func New(svc *service.Service) *warden.Server {
	s := warden.NewServer(nil)
	s.Use(middleware())
	_, err := s.Start()
	if err != nil {
		panic("run server failed!" + err.Error())
	}
	return s
}

func middleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//call chain
		resp, err = handler(ctx, req)
		return
	}
}
