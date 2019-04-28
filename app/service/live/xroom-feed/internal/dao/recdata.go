package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/service/live/xroom-feed/internal/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

var (
	recConfKey  = "rec_conf"
	recPoolKey  = "rec_pool_%d"
	recRoomInfo = "rec_info_%d"
)

//GetCacheData 获取缓存数据
func (d *Dao) GetCacheData() (data []byte) {
	conn := d.redis.Get(context.Background())
	defer conn.Close()

	var err error
	if data, err = redis.Bytes(conn.Do("GET", recConfKey)); err != nil {
		log.Error("[recdata cache] GetCacheData err:%+v", err)
		return make([]byte, 0)
	}
	return data
}

//GetRecInfoByRoomid 批量通过Roomid获取rec_info
func (d *Dao) GetRecInfoByRoomid(ctx context.Context, roomids []int64) map[int64]*model.RecRoomInfo {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	redisSuccess := make([]int64, 0, len(roomids))
	for _, roomid := range roomids {
		key := fmt.Sprintf(recRoomInfo, roomid)
		if err := conn.Send("HGETALL", key); err != nil {
			log.Error("[recdata] GetRecInfoByRoomid redis send roomid:%d err:%v", roomid, err)
			continue
		}
		redisSuccess = append(redisSuccess, roomid)
	}

	conn.Flush()

	result := make(map[int64]*model.RecRoomInfo)
	for _, roomid := range redisSuccess {
		var roominfo map[string]string
		var err error
		if roominfo, err = redis.StringMap(conn.Receive()); err != nil {
			log.Error("[recdata] GetRecInfoByRoomid redis receive err:%d err:%+v", roomid, err)
			continue
		}
		if len(roominfo) <= 0 {
			log.Error("[recdata] GetRecInfoByRoomid redis receive empty:%d err:+%v", roomid, roominfo)
			continue
		}

		result[roomid] = model.NewRecRoomInfo()
		result[roomid].Title = roominfo["title"]
		result[roomid].Uid, _ = strconv.ParseInt(roominfo["uid"], 10, 64)
		result[roomid].PopularityCount, _ = strconv.ParseInt(roominfo["popularity_count"], 10, 64)
		result[roomid].KeyFrame = roominfo["Keyframe"]
		result[roomid].Cover = roominfo["cover"]
		result[roomid].ParentAreaID, _ = strconv.ParseInt(roominfo["parent_area_id"], 10, 64)
		result[roomid].ParentAreaName = roominfo["parent_area_name"]
		result[roomid].AreaID, _ = strconv.ParseInt(roominfo["area_id"], 10, 64)
		result[roomid].AreaName = roominfo["area_name"]
	}
	return result
}

//GetRecPoolByID  通过id获取推荐池
func (d *Dao) GetRecPoolByID(ctx context.Context, rids []int64) map[int64][]int64 {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	var (
		room string
		err  error
	)

	redisSuccess := make([]int64, 0, len(rids))
	for _, rid := range rids {
		key := fmt.Sprintf(recPoolKey, rid)
		if err := conn.Send("GET", key); err != nil {
			log.Error("[recdata] GetRecPoolByID send id %d err:%+v", rid, err)
			continue
		}
		redisSuccess = append(redisSuccess, rid)
	}

	conn.Flush()

	result := make(map[int64][]int64)
	for _, rid := range redisSuccess {
		if room, err = redis.String(conn.Receive()); err != nil {
			log.Error("[recdata] GetRecPoolByID receive id %d err:%+v", rid, err)
			continue
		}

		if room == "" {
			log.Info("[recdata] GetRecPoolByID receive empty room list rid: %d", rid)
			continue
		}
		rooms := strings.Split(room, ",")
		roomIDs := make([]int64, len(rooms), len(rooms))

		for i, roomid := range rooms {
			roomIDs[i], _ = strconv.ParseInt(roomid, 10, 64)
		}
		result[rid] = roomIDs
	}

	return result
}
