package v1

import (
	"context"
	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"
)

// SplashService struct
type SplashService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

//NewSplashService init
func NewSplashService(c *conf.Config) (s *SplashService) {
	s = &SplashService{
		conf: c,
	}
	return s
}

// GetInfo implementation
// 获取有效闪屏配置
// `method:"GET" internal:"true" `
func (s *SplashService) GetInfo(ctx context.Context, req *v1pb.GetInfoReq) (resp *v1pb.GetInfoResp, err error) {
	resp = &v1pb.GetInfoResp{}
	return
}
