package grpc

import (
	"fmt"
	pb "go-common/app/service/live/gift/api/grpc/v1"
	"go-common/app/service/live/gift/conf"
	"go-common/app/service/live/gift/dao"
	svc "go-common/app/service/live/gift/service/v1"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

//Init Init
func Init(c *conf.Config) {
	s := warden.NewServer(nil) // 酌情传入config
	dao.InitApi()
	gs := svc.NewGiftService(c)
	pb.RegisterGiftServer(s.Server(), gs)
	_, err := s.Start()
	if err != nil {
		log.Error("grpc Start error(%v)", err)
		panic(err)
	}
	fmt.Println("start")
	gs.TickerReloadGift()
}
