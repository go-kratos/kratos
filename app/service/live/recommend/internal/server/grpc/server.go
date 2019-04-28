package grpc

import (
	v1pb "go-common/app/service/live/recommend/api/grpc/v1"
	"go-common/app/service/live/recommend/internal/conf"
	svc "go-common/app/service/live/recommend/internal/service/v1"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

// Init grpc server
func Init(c *conf.Config) {
	s := warden.NewServer(nil)
	v1pb.RegisterRecommendServer(s.Server(), svc.NewRecommendService(c))
	_, err := s.Start()
	if err != nil {
		log.Error("grpc Start error(%v)", err)
		panic(err)
	}
}
