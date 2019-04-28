package v1

import (
	"context"
	"encoding/json"

	v1pb "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/app/service/live/resource/conf"
	"go-common/app/service/live/resource/dao"
	"go-common/library/ecode"
)

// SplashService struct
type SplashService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewSplashService init
func NewSplashService(c *conf.Config) (s *SplashService) {
	s = &SplashService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

const typeSplash = "splash"

// GetInfo implementation
// 获取有效闪屏配置
func (s *SplashService) GetInfo(ctx context.Context, req *v1pb.GetInfoReq) (resp *v1pb.GetInfoResp, err error) {
	resp = &v1pb.GetInfoResp{}
	reply, err := s.dao.GetInfo(ctx, typeSplash, req.Platform, req.Build)
	if err != nil {
		err = ecode.GetSplashErr
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
