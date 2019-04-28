package dao

import (
	"context"

	"go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/log"
)

//PERMANENT_LOCK_TIME 永久封禁时间
const PERMANENT_LOCK_TIME = "2037-01-01 00:00:00"

//GetRoomInfosByUids ...
func (d *Dao) GetRoomInfosByUids(c context.Context, uid []int64) (res map[int64]*v1.RoomGetStatusInfoByUidsResp_RoomInfo, err error) {
	reply, err := d.RoomApi.V1Room.GetStatusInfoByUids(c, &v1.RoomGetStatusInfoByUidsReq{
		Uids:              uid,
		FilterOffline:     0,
		ShowHidden:        1,
		FilterIndexBlack:  0,
		FilterVideo:       0,
		NeedBroadcastType: 0,
	})
	if err != nil {
		log.Error("room_getRoomInfosByUids_error:%v", err)
		return
	}
	if reply.Code != 0 {
		log.Error("room_getRoomInfosByUids_error:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
		return
	}
	log.Info("room_getRoomInfosByUids:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
	res = reply.Data
	return
}

//CreateRoom ...
func (d *Dao) CreateRoom(c context.Context, uid int64) (res *v1.RoomMngCreateRoomResp_Data, err error) {
	reply, err := d.RoomApi.V1RoomMng.CreateRoom(c, &v1.RoomMngCreateRoomReq{Uid: uid})
	if err != nil {
		log.Error("room_createRoom_error:%v", err)
		return
	}
	if reply.Code != 0 {
		log.Error("room_createRoom_error:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
		return
	}
	log.Info("room_createRoom:%d,%s,$v", reply.Code, reply.Msg, reply.Data)
	res = reply.Data
	return
}
