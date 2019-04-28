package room

import (
	"context"
	"errors"
	"time"

	"go-common/app/interface/live/app-interface/conf"
	cDao "go-common/app/interface/live/app-interface/dao"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
)

func (d *Dao) GetActivityCard(ctx context.Context, ids []int64, logPrefix string) (err error, data map[int64]*roomV2.AppIndexGetActivityCardResp_ActivityCard) {
	activityQueryTimeout := time.Duration(conf.GetTimeout("activityQuery", 100)) * time.Millisecond

	cardInfo, roomError := cDao.RoomApi.V2AppIndex.GetActivityCard(rpcCtx.WithTimeout(ctx, activityQueryTimeout), &roomV2.AppIndexGetActivityCardReq{Ids: ids})

	if roomError != nil {
		log.Error("[%s] get activity card info rpc error, room.v2.AppIndex.GetActivityCard, error:%+v", logPrefix, roomError)
		err = errors.New("get activity card info rpc error")
		return
	}
	if cardInfo.Code != 0 {
		log.Error("[%s] get activity card info response error, code, %d, msg: %s", logPrefix, cardInfo.Code, cardInfo.Msg)
		err = errors.New("get activity card info response error")
		return
	}
	if cardInfo.Data.ActivityCard == nil {
		log.Error("[%s] get activity card info but on data", logPrefix)
		err = errors.New("get activity card info but no data")
		return
	}
	data = cardInfo.Data.ActivityCard
	return
}
