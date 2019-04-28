package v1

import (
	resAPI "go-common/app/interface/live/app-room/api/http/v1"
	"go-common/app/interface/live/app-room/conf"
	rspb "go-common/app/service/live/resource/api/grpc/v1"
	bm "go-common/library/net/http/blademaster"
)

var _rsCli *rspb.Client

// Init -
func Init(c *conf.Config) {
	var err error
	if _rsCli, err = rspb.NewClient(c.ResourceClient); err != nil {
		panic(err)
	}
}

// GetBanner -
func GetBanner(context *bm.Context) {
	p := new(rspb.GetBannerReq)
	if err := context.Bind(p); err != nil {
		return
	}
	respRPC, err := _rsCli.GetBanner(context, p)
	if err != nil {
		return
	}
	resp := make([]resAPI.GetBannerResp, len(respRPC.List))
	for index, banner := range respRPC.List {
		resp[index].Id = banner.Id
		resp[index].Title = banner.Title
		resp[index].Img = banner.ImageUrl
		resp[index].Link = banner.JumpPath
	}
	context.JSON(resp, err)
}
