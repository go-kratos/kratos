package grpc

import (
	pb "go-common/app/service/live/grpc-demo/api/grpc/v1"
	"go-common/app/service/live/grpc-demo/conf"
	svc "go-common/app/service/live/grpc-demo/service/v1"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

// TODO

func Init(c *conf.Config) {
	s := warden.NewServer(nil)
	pb.RegisterGreeterServer(s.Server(), svc.NewGreeterService(c))
	_, err := s.Start()
	if err != nil {
		log.Error("grpc Start error(%v)", err)
		panic(err)
	}
}
