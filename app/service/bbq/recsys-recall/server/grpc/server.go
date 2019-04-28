package grpc

import (
	"context"
	"flag"
	"time"

	v1 "go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	"go-common/app/service/bbq/recsys-recall/service"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"

	"google.golang.org/grpc"
)

var (
	_gRPCAddr string
)

func init() {
	flag.StringVar(&_gRPCAddr, "grpc_addr", "0.0.0.0:9000", "default config path")
}

//New 生成rpc服务
func New(srv *service.Service) *warden.Server {
	servConf := &warden.ServerConfig{
		Addr:    _gRPCAddr,
		Timeout: xtime.Duration(2 * time.Second),
	}
	s := warden.NewServer(servConf)
	s.Use(middleware())
	v1.RegisterRecsysRecallServer(s.Server(), srv)
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
