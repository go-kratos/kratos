package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/library/log"
)

type redisKeyResp struct {
	RedisKey string `json:"redis_key"`
	TimeOut  int    `json:"time_out"`
}

//RoomRecordListForLive 房间侧开播记录相关历史数据，存3天
func (d *Dao) LRoomLiveRecordList(content string, roomId int64, liveTime int64) (resp *redisKeyResp) {
	resp = &redisKeyResp{}
	if liveTime <= 0 {
		roomInfo, err := d.FetchRoomByIDs(context.TODO(), &v1.RoomByIDsReq{RoomIds: []int64{roomId}, Fields: []string{"live_start_time"}})
		if err != nil || roomInfo.RoomDataSet == nil {
			log.Error("LRoomLiveRecordList_err:err=%v;info=%v", err, roomInfo)
			return
		}
		liveTime = roomInfo.RoomDataSet[roomId].LiveStartTime
	}
	contentType := content + "_" + strconv.Itoa(int(liveTime))
	resp.RedisKey = fmt.Sprintf(contentType+"_list_%d", roomId)
	resp.TimeOut = 24 * 60 * 60 * 3
	return
}

//RoomRecordList 房间侧记录相关历史数据，存1天
func (d *Dao) LRoomRecordList(contentType string, roomId int64) (resp *redisKeyResp) {
	resp = &redisKeyResp{}
	resp.RedisKey = fmt.Sprintf(contentType+"_list_%d", roomId)
	resp.TimeOut = 24 * 60 * 60
	return
}

//RoomRecordCurrent 房间侧实时数据记录，存一个小时
func (d *Dao) SRoomRecordCurrent(content string, roomId int64) (resp *redisKeyResp) {
	resp = &redisKeyResp{}
	resp.RedisKey = fmt.Sprintf(content+"_key_%d", roomId)
	resp.TimeOut = 60 * 60
	return
}

//主播侧开播记录相关历史数据，存3天
func (d *Dao) LUserLiveRecordList(content string, uid int64, liveTime int64) (resp *redisKeyResp) {
	resp = &redisKeyResp{}
	if liveTime <= 0 {
		roomInfo, err := d.FetchRoomByIDs(context.TODO(), &v1.RoomByIDsReq{Uids: []int64{uid}, Fields: []string{"live_start_time"}})
		if err != nil {
			log.Error("LRoomLiveRecordList_err:err=%v;info=%v", err, roomInfo)
			return
		}
	}
	contentType := content + "_" + strconv.Itoa(int(liveTime))
	resp.RedisKey = fmt.Sprintf(contentType+"_list_%d", uid)
	resp.TimeOut = 24 * 60 * 60 * 3
	return
}

//主播侧记录相关历史数据，存1天
func (d *Dao) LUserRecordList(contentType string, uid int64) (resp *redisKeyResp) {
	resp = &redisKeyResp{}
	resp.RedisKey = fmt.Sprintf(contentType+"_list_%d", uid)
	resp.TimeOut = 24 * 60 * 60
	return
}

//UserRecordCurrent 主播侧实时数据记录，存一个小时
func (d *Dao) SUserRecordCurrent(content string, uid int64) (resp *redisKeyResp) {
	resp = &redisKeyResp{}
	resp.RedisKey = fmt.Sprintf(content+"_key_%d", uid)
	resp.TimeOut = 60 * 60
	return
}
