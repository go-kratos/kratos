package room_ex

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-common/app/interface/live/app-interface/conf"
	cDao "go-common/app/interface/live/app-interface/dao"
	roomEx "go-common/app/service/live/room_ex/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
	"time"
)

func (d *Dao) GetBanner(ctx context.Context, biz int64, position int64, platform string, device string, build int64) (bannerListResp []*roomEx.BannerGetNewBannerResp_NewBanner, err error) {
	timeOut := time.Duration(conf.GetTimeout("/room_ex/v1/Banner/getNewBanner", 50)) * time.Millisecond
	bannerListResult, errCode := cDao.RoomExtApi.V1Banner.GetNewBanner(rpcCtx.WithTimeout(ctx, timeOut), &roomEx.BannerGetNewBannerReq{
		Platform: biz,
		Position: position,
		UserPlatform: platform,
		UserDevice: device,
		Build: build,
	})

	if errCode != nil {
		log.Error("[getBannerFromRoomEx]roomEx.v1.getNewBanner rpc error:%+v", errCode)
		err = errors.WithMessage(ecode.GetBannerErr, fmt.Sprintf("roomEx.v1.getNewBanner rpc error:%+v", errCode))
		return
	}
	if bannerListResult.Code != 0 {
		log.Error("[getBannerFromRoomEx]roomEx.v1.getNewBanner response error:%+v,code:%d,msg:%s", errCode, bannerListResult.Code, bannerListResult.Msg)
		err = errors.WithMessage(ecode.GetBannerErr, fmt.Sprintf("roomEx.v1.getNewBanner response error,code:%d,msg:%s", bannerListResult.Code, bannerListResult.Msg))
		return
	}

	if bannerListResult.Data == nil {
		log.Error("[getBannerFromRoomEx]roomEx.v1.getNewBanner empty error")
		err = errors.WithMessage(ecode.GetBannerErr, "roomEx.v1.getNewBanner empty error")
		return
	}
	bannerListResp = bannerListResult.Data

	return
}