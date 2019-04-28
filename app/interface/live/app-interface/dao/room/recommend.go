package room

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"go-common/app/interface/live/app-interface/conf"
	cDao "go-common/app/interface/live/app-interface/dao"
	roomV1 "go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
)

// 获取首页推荐强推列表
func (d *Dao) GetStrongRecList(ctx context.Context, recPage int64) (strongRecRoomListResp *roomV1.RoomRecommendClientRecStrongResp, err error) {
	strongRecRoomListResp = &roomV1.RoomRecommendClientRecStrongResp{}
	clientRecStrongTimeout := time.Duration(conf.GetTimeout("clientRecStrong", 100)) * time.Millisecond
	strongRecRoomList, err := cDao.RoomApi.V1RoomRecommend.ClientRecStrong(rpcCtx.WithTimeout(ctx, clientRecStrongTimeout), &roomV1.RoomRecommendClientRecStrongReq{
		RecPage: recPage,
	})
	if err != nil {
		log.Error("[GetStrongRecList]room.v1.clientStrongRec rpc error:%+v", err)
		err = errors.New(fmt.Sprintf("room.v1.clientStrongRec rpc error:%+v", err))
		return
	}
	if strongRecRoomList.Code != 0 {
		log.Error("[GetStrongRecList]room.v1.getPendantByIds response code:%d,msg:%s", strongRecRoomList.Code, strongRecRoomList.Msg)
		err = errors.New(fmt.Sprintf("room.v1.getPendantByIds response error, code:%d, msg:%s", strongRecRoomList.Code, strongRecRoomList.Msg))
		return
	}

	if strongRecRoomList.Data == nil || strongRecRoomList.Data.Result == nil {
		log.Error("[GetStrongRecList]room.v1.getPendantByIds empty")
		err = errors.New("[getSkyHorseRoomList]room.v1.getPendantByIds empty")
		return
	}

	strongRecRoomListResp = strongRecRoomList
	return
}
