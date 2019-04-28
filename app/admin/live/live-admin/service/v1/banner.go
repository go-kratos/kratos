package v1

import (
	"context"
	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"
)

// BannerService struct
type BannerService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

//NewBannerService init
func NewBannerService(c *conf.Config) (s *BannerService) {
	s = &BannerService{
		conf: c,
	}
	return s
}

// GetBlinkBanner implementation
// 获取有效banner配置
// `method:"GET" internal:"true" `
func (s *BannerService) GetBlinkBanner(ctx context.Context, req *v1pb.GetInfoReq) (resp *v1pb.GetInfoResp, err error) {
	resp = &v1pb.GetInfoResp{}
	return
}

// GetBanner implementation
// 获取有效banner配置
// `method:"GET" internal:"true" `
func (s *BannerService) GetBanner(ctx context.Context, req *v1pb.GetBannerReq) (resp *v1pb.GetBannerResp, err error) {
	resp = &v1pb.GetBannerResp{}
	return
}
