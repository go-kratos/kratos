package room

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"go-common/app/interface/live/app-interface/conf"
	cDao "go-common/app/interface/live/app-interface/dao"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
)

// GetRoomInfoByIds
func (d *Dao) GetRoomInfoByIds(ctx context.Context, roomIds []int64, fields []string, from string) (multiRoomListResp map[int64]*roomV2.RoomGetByIdsResp_RoomInfo, err error) {
	multiRoomListResp = make(map[int64]*roomV2.RoomGetByIdsResp_RoomInfo)
	getByIdsTimeout := time.Duration(conf.GetTimeout("getByIds", 100)) * time.Millisecond
	multiRoomList, getByIdsError := cDao.RoomApi.V2Room.GetByIds(rpcCtx.WithTimeout(ctx, getByIdsTimeout), &roomV2.RoomGetByIdsReq{
		Ids:               roomIds,
		NeedBroadcastType: 1,
		NeedUinfo:         1,
		Fields:            fields,
		From:              from,
	})

	if getByIdsError != nil {
		log.Error("[GetRoomInfoByIds]room.v2.getByIds rpc error:%+v", getByIdsError)
		// 这个是推荐房间列表的基础信息，如果失败需要cancel，不然返回值会很奇怪
		err = errors.WithMessage(ecode.GetRoomError, fmt.Sprintf("room.v2.getByIds rpc error:%+v", getByIdsError))
		return
	}
	if multiRoomList.Code != 0 {
		log.Error("[GetRoomInfoByIds]room.v2.getByIds response error,code:%d,msg:%s", multiRoomList.Code, multiRoomList.Msg)
		// 这个是推荐房间列表的基础信息，如果失败需要cancel，不然返回值会很奇怪
		err = errors.WithMessage(ecode.GetRoomError, fmt.Sprintf("room.v2.getByIds response error,code:%d,msg:%s", multiRoomList.Code, multiRoomList.Msg))
		return
	}

	if multiRoomList.Data == nil {
		log.Error("[GetRoomInfoByIds]room.v2.getByIds empty error")
		// 这个是推荐房间列表的基础信息，如果失败需要cancel，不然返回值会很奇怪
		err = errors.WithMessage(ecode.GetRoomEmptyError, "room.v2.getByIds empty error")
		return
	}
	multiRoomListResp = multiRoomList.Data

	return
}
