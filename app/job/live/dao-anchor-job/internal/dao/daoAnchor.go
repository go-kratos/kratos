package dao

import (
	"context"
	"time"

	"go-common/library/ecode"

	"fmt"

	daoAnchorV0 "go-common/app/service/live/dao-anchor/api/grpc/v0"
	"go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/library/log"
)

//LIVE_OPEN 开播
const LIVE_OPEN = 1

//LIVE_CLOSE 关播
const LIVE_CLOSE = 0

//LIVE_ROUND 轮播
const LIVE_ROUND = 2

//PAGE_SIZE 分页数据量
const PAGE_SIZE = 100

//RETRY_TIME  接口充实次数
const RETRY_TIME = 3

//GetAllLiveRoom 获取在播列表
func (d *Dao) GetAllLiveRoom(c context.Context, fields []string) (resp map[int64]*v1.RoomData, err error) {
	page := 0
	resp = map[int64]*v1.RoomData{}
	retry := 1
	for true {
		reply, err := d.daoAnchorApi.RoomOnlineList(c, &v1.RoomOnlineListReq{Fields: fields, Page: int64(page), PageSize: PAGE_SIZE})
		if err != nil {
			if retry >= RETRY_TIME {
				break
			}
			retry++
			time.Sleep(time.Second * 3)
			log.Errorw(c, "log", fmt.Sprintf("getAllLiveRoom_RoomOnlineList_error:page=%d;err=%v", page, err))
			continue
		}
		if len(reply.RoomDataList) <= 0 {
			break
		}
		page = page + 1
		roomDataList := reply.RoomDataList
		for roomId, info := range roomDataList {
			v := info
			resp[roomId] = v
		}
		time.Sleep(time.Millisecond)
	}
	return
}

//GetAllLiveRoomIds 获取在播列表
func (d *Dao) GetAllLiveRoomIds(c context.Context) (resp []int64, err error) {
	page := 0
	reply, err := d.daoAnchorApi.RoomOnlineListByArea(c, &v1.RoomOnlineListByAreaReq{})
	if err != nil {
		log.Errorw(c, "log", fmt.Sprintf("GetAllLiveRoomIds_RoomOnlineList_error:page=%d;err=%v", page, err))
		return
	}
	if len(reply.RoomIds) <= 0 {
		return
	}
	resp = reply.RoomIds
	return
}

//GetInfosByRoomIds  获取主播房间信息
func (d *Dao) GetInfosByRoomIds(c context.Context, roomIds []int64, fields []string) (resp map[int64]*v1.RoomData, err error) {
	if roomIds == nil {
		err = ecode.InvalidParam
		log.Errorw(c, "log", fmt.Sprintf("getInfosByRoomIds_params_error:%v", roomIds))
		return
	}
	reply, err := d.daoAnchorApi.FetchRoomByIDs(c, &v1.RoomByIDsReq{RoomIds: roomIds, Fields: fields})
	if err != nil {
		log.Errorw(c, "log", fmt.Sprintf("getInfosByRoomIds_FetchRoomByIDs_error:reply=%v;err=%v", reply, err))
		return
	}
	if reply == nil {
		err = ecode.CallDaoAnchorError
		log.Errorw(c, "log", "getInfosByRoomIds_FetchRoomByIDs_error")
		return
	}
	resp = reply.RoomDataSet
	return
}

//UpdateRoomEx ...
func (d *Dao) UpdateRoomEx(c context.Context, roomId int64, fields []string, keyFrame string) (resp int64, err error) {
	if roomId < 0 {
		err = ecode.InvalidParam
		log.Errorw(c, "log", fmt.Sprintf("updateRoom_params_error:%v", roomId))
		return
	}
	reply, err := d.daoAnchorApi.RoomExtendUpdate(c, &v1.RoomExtendUpdateReq{Fields: fields, RoomId: roomId, Keyframe: keyFrame})
	if err != nil {
		log.Errorw(c, "log", fmt.Sprintf("updateRoom_RoomUpdate_error:reply=%v;err=%v", reply, err))
		return
	}
	resp = reply.AffectedRows
	return
}

func (d *Dao) CreateCacheList(c context.Context, roomIds []int64, content string) (err error) {
	if len(roomIds) <= 0 || content == "" {
		log.Errorw(c, "log", fmt.Sprintf("CreateCacheList_params_error:room_id=%v;content=%s", roomIds, content))
		return
	}
	reply, err := d.daoAnchorApiV0.CreateLiveCacheList(c, &daoAnchorV0.CreateLiveCacheListReq{RoomIds: roomIds, Content: content})
	if err != nil {
		log.Errorw(c, "log", fmt.Sprintf("CreateCacheList_error:reply=%v;err=%v", reply, err))
		return
	}
	log.Info("CreateCacheList_info:roomIds=%v;content=%s;reply=%v", roomIds, content, reply)
	return
}
func (d *Dao) GetContentMap(c context.Context) (resp map[string]int64, err error) {
	reply, err := d.daoAnchorApiV0.GetContentMap(c, &daoAnchorV0.GetContentMapReq{})
	if err != nil {
		log.Errorw(c, "log", fmt.Sprintf("GetContentMap_error:reply=%v;err=%v", reply, err))
		return
	}
	resp = reply.List
	return
}
func (d *Dao) CreateDBData(c context.Context, roomIds []int64, content string) (err error) {
	if len(roomIds) <= 0 || content == "" {
		log.Errorw(c, "log", fmt.Sprintf("CreateCacheList_params_error:room_id=%v;content=%s", roomIds, content))
		return
	}
	reply, err := d.daoAnchorApiV0.CreateDBData(c, &daoAnchorV0.CreateDBDataReq{RoomIds: roomIds, Content: content})
	if err != nil {
		log.Errorw(c, "log", fmt.Sprintf("CreateCacheList_error:reply=%v;err=%v", reply, err))
		return
	}
	return
}
