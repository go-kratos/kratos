package grpc

import (
	pb "go-common/app/service/live/xcaptcha/api/grpc/v1"
	"go-common/app/service/live/xcaptcha/conf"
	svc "go-common/app/service/live/xcaptcha/service/v1"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

// Init gRpc Init
func Init(c *conf.Config) {
	s := warden.NewServer(nil) // 酌情传入config
	// 每个proto里定义的service添加一行
	pb.RegisterXCaptchaServer(s.Server(), svc.NewXCaptchaService(c))
	_, err := s.Start()
	if err != nil {
		log.Error("grpc Start error(%v)", err)
		panic(err)
	}
}
