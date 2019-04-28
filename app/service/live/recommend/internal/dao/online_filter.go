package dao

import (
	"context"
	"sync/atomic"
	"time"

	"go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/log"
)

var onlineRoomIdValue atomic.Value
var areaRoomsValue atomic.Value

// StartRefreshJob 更新在线房间信息
func StartRefreshJob() {
	t := time.Tick(time.Second * 60)
	refreshOnlineRoomData(context.Background())
	for range t {
		refreshOnlineRoomData(context.Background())
	}
}

// refreshOnlineRoomData 更新RoomId
func refreshOnlineRoomData(ctx context.Context) (err error) {
	resp, err := RoomAPI.V1Room.AllLiveForBigdata(ctx, &v1.RoomAllLiveForBigdataReq{})
	if err != nil {
		return
	}
	onlineRooms := map[int64]struct{}{}
	areaRooms := map[int64][]int64{}
	for _, info := range resp.Data {
		roomID := info.Roomid

		onlineRooms[roomID] = struct{}{}
		if info.Online > 100 {
			areaRooms[info.AreaV2Id] = append(areaRooms[info.AreaV2Id], roomID)
		}
	}
	log.Info("refreshOnlineRoomData: count=%d", len(onlineRooms))
	log.Info("refreshOnlineRoomData area Rooms: %+v", areaRooms)
	onlineRoomIdValue.Store(onlineRooms)
	areaRoomsValue.Store(areaRooms)
	return
}

func (d *Dao) getAreaRoomIds(areaId int64) (ret []int64) {
	ret = make([]int64, 0)
	areaRooms, ok := areaRoomsValue.Load().(map[int64][]int64)
	if !ok {
		log.Warn("cannot load current area room ids")
		return
	}
	ret = areaRooms[areaId]
	return
}

// FilterOnlineRoomIds 给定一批room id 返回所有在线的
func (d *Dao) FilterOnlineRoomIds(roomIds []int64) (ret []int64) {
	ret = make([]int64, 0)
	currentIds, ok := onlineRoomIdValue.Load().(map[int64]struct{})
	if !ok {
		log.Warn("cannot load current online room ids")
		return
	}
	for _, roomId := range roomIds {
		if _, ok := currentIds[roomId]; ok {
			ret = append(ret, roomId)
		}
	}
	return
}
