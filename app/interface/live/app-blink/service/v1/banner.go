package v1

import (
	"context"

	"go-common/app/interface/live/app-blink/conf"
	rspb "go-common/app/service/live/resource/api/grpc/v1"
)

// BannerService struct
type BannerService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	rsCli *rspb.Client
}

//NewBannerService init
func NewBannerService(c *conf.Config) (s *BannerService) {
	s = &BannerService{
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

// GetBlinkBanner implementation
// 获取banner配置
// `dynamic:"true"`
func (s *BannerService) GetBlinkBanner(ctx context.Context, req *rspb.GetInfoReq) (resp *rspb.GetInfoResp, err error) {
	resp, err = s.rsCli.GetBlinkBanner(ctx, req)
	return
}
