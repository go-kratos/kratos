package v1

import (
	"context"
	"encoding/json"

	v1pb "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/app/service/live/resource/conf"
	"go-common/app/service/live/resource/dao"
	"go-common/library/ecode"
)

// BannerService struct
type BannerService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewBannerService init
func NewBannerService(c *conf.Config) (s *BannerService) {
	s = &BannerService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

const typeBanner = "banner"

// GetBlinkBanner implementation
// 获取有效闪屏配置
func (s *BannerService) GetBlinkBanner(ctx context.Context, req *v1pb.GetInfoReq) (resp *v1pb.GetInfoResp, err error) {
	resp = &v1pb.GetInfoResp{}
	reply, err := s.dao.GetInfo(ctx, typeBanner, req.Platform, req.Build)
	if err != nil {
		err = ecode.GetBannerErr
		return
	}
	if reply == nil || reply.ID < 1 {
		err = ecode.NothingFound
		return
	}
	type updateImage struct {
		JumpPath     string `json:"JumpPath"`
		JumpPathType int64  `json:"JumpPathType"`
		JumpTime     int64  `json:"JumpTime"`
		ImageUrl     string `json:"ImageUrl"`
	}
	imageInfo := reply.ImageInfo
	imageInfoArr := &updateImage{}
	json.Unmarshal([]byte(imageInfo), imageInfoArr)

	resp.Id = reply.ID
	resp.Title = reply.Title
	resp.ImageUrl = imageInfoArr.ImageUrl
	resp.JumpTime = imageInfoArr.JumpTime
	resp.JumpPath = imageInfoArr.JumpPath
	resp.JumpPathType = imageInfoArr.JumpPathType
	return
}

// GetBanner implementation
func (s *BannerService) GetBanner(ctx context.Context, req *v1pb.GetBannerReq) (resp *v1pb.GetBannerResp, err error) {
	resp = &v1pb.GetBannerResp{}
	banners, err := s.dao.GetBanner(ctx, req.Platform, req.Build, req.Type)
	if err != nil {
		err = ecode.GetBannerErr
		return
	}
	type updateImage struct {
		JumpPath     string `json:"JumpPath"`
		JumpPathType int64  `json:"JumpPathType"`
		JumpTime     int64  `json:"JumpTime"`
		ImageUrl     string `json:"ImageUrl"`
	}
	if banners == nil {
		return
	}
	for _, banner := range banners {
		if banner.ID < 1 {
			continue
		}
		imageInfo := banner.ImageInfo
		imageInfoArr := &updateImage{}
		json.Unmarshal([]byte(imageInfo), imageInfoArr)
		b := &v1pb.GetBannerResp_List{}
		b.Id = banner.ID
		b.Title = banner.Title
		b.ImageUrl = imageInfoArr.ImageUrl
		b.JumpTime = imageInfoArr.JumpTime
		b.JumpPath = imageInfoArr.JumpPath
		b.JumpPathType = imageInfoArr.JumpPathType
		resp.List = append(resp.List, b)
	}
	return
}
