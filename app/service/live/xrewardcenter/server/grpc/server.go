package grpc

import (
	pb "go-common/app/service/live/xrewardcenter/api/grpc/v1"
	"go-common/app/service/live/xrewardcenter/conf"
	"go-common/app/service/live/xrewardcenter/dao"
	svc "go-common/app/service/live/xrewardcenter/service/v1"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

// TODO

// Init .
func Init(c *conf.Config) {
	s := warden.NewServer(nil)
	dao.InitAPI()
	pb.RegisterAnchorRewardServer(s.Server(), svc.NewAnchorTaskService(c))
	_, err := s.Start()
	if err != nil {
		log.Error("grpc Start error(%v)", err)
		panic(err)
	}
}
