package dao

import (
	"context"

	daoAnchorV1 "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/library/log"
)

func (d *Dao) SetRecPoolCache(ctx context.Context, key string, list string, expire int) (err error) {
	log.Info("[SetRecPoolCache]key: %s", key)
	rconn := d.redis.Get(ctx)
	defer rconn.Close()
	if _, err = rconn.Do("SET", key, list, "EX", expire); err != nil {
		log.Error("[SetRecPoolCache]redis.Set error(%v)", err)
	}
	return
}

func (d *Dao) SetRecInfoCache(ctx context.Context, list map[int64]*daoAnchorV1.RoomData) (err error) {
	rconn := d.redis.Get(ctx)
	defer rconn.Close()
	//rconn.Send("MULTI")
	for roomId, roomData := range list {
		if roomData == nil {
			continue
		}
		key, expire := d.getRecInfoKey(roomId)
		rconn.Send("HMSET", key, "room_id", roomData.RoomId, "uid", roomData.Uid, "title", roomData.Title, "popularity_count", roomData.PopularityCount, "Keyframe", roomData.Keyframe, "cover", roomData.Cover, "parent_area_id", roomData.ParentAreaId, "parent_area_name", roomData.ParentAreaName, "area_id", roomData.AreaId, "area_name", roomData.AreaName)
		rconn.Send("EXPIRE", key, expire)
	}
	if err = rconn.Flush(); err != nil {
		log.Error("[SetRecInfoCache]redis.Set error(%v)", err)
	}

	for _, roomData := range list {
		if roomData == nil {
			continue
		}
		rconn.Receive()
		rconn.Receive()
	}

	return
}
