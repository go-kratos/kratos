package dao

import (
	"context"
	"time"

	"go-common/app/service/video/stream-mng/api/v1"
	"go-common/library/log"
)

//视频云接口调用

//GetPicsByRoomId 根据房间id获取当前关键帧
func (d *Dao) GetPicsByRoomId(c context.Context, roomId int64, startTime time.Time, endTime time.Time) (resp []string, err error) {
	reply, err := d.VideoApi.GetSingleScreeShot(c, &v1.GetSingleScreeShotReq{RoomId: roomId, StartTime: startTime.Format("2006-01-02 15:04:05"), EndTime: endTime.Format("2006-01-02 15:04:05")})
	if err != nil {
		log.Error("getPicsByRoomId_GetSingleScreeShot_error:reply:%v;err=%v", reply, err)
		return
	}
	if reply == nil {
		log.Error("getPicsByRoomId_GetSingleScreeShot_error")
		return
	}
	resp = reply.List
	return
}
