package v1

import (
	"context"

	"go-common/app/interface/live/app-blink/conf"
	rspb "go-common/app/service/live/resource/api/grpc/v1"
)

// SplashService struct
type SplashService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	rsCli *rspb.Client
}

//NewSplashService init
func NewSplashService(c *conf.Config) (s *SplashService) {
	s = &SplashService{
		conf: c,
	}
	var svc *rspb.Client
	var err error

	if svc, err = rspb.NewClient(s.conf.ResourceClient); err != nil {
		panic(err)
	}
	s.rsCli = svc
	return s
}

// GetInfo implementation
// 获取有效闪屏配置
// `dynamic:"true"`
func (s *SplashService) GetInfo(ctx context.Context, req *rspb.GetInfoReq) (resp *rspb.GetInfoResp, err error) {
	resp, err = s.rsCli.GetInfo(ctx, req)
	return
}
