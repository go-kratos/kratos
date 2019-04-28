package grpc

import (
	"fmt"

	"go-common/app/service/live/xlottery/dao"

	pb "go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/app/service/live/xlottery/conf"
	svc "go-common/app/service/live/xlottery/service/v1"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

// Init .
func Init(c *conf.Config) *warden.Server {
	dao.InitAPI()
	s := warden.NewServer(nil) // 酌情传入config
	gs := svc.NewCapsuleService(c)
	pb.RegisterCapsuleServer(s.Server(), gs)
	pb.RegisterStormServer(s.Server(), svc.NewStromService(c))
	_, err := s.Start()
	if err != nil {
		log.Error("grpc Start error(%v)", err)
		panic(err)
	}
	fmt.Println("start")
	return s
}
