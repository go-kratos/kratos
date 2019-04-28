package roomAdmin

import (
	"context"
	"fmt"
	"go-common/app/service/live/xuser/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

// AddCacheNoneUser write an flag in cache represents empty
func (d *Dao) AddCacheNoneUser(c context.Context, uid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := KeyUser(uid)

	var roomAdmins []model.RoomAdmin
	roomAdmins = append(roomAdmins, model.RoomAdmin{
		Id: -1,
	})

	item := &memcache.Item{Key: key, Object: roomAdmins, Expiration: d.RoomAdminExpire, Flags: memcache.FlagJSON | memcache.FlagGzip}
	if err = conn.Set(item); err != nil {
		prom.BusinessErrCount.Incr("mc:AddCacheRoomAdminAnchor")
		log.Errorv(c, log.KV("AddCacheRoomAdminAnchor", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// AddCacheNoneRoom write an flag in cache represents empty
func (d *Dao) AddCacheNoneRoom(c context.Context, uid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := KeyRoom(uid)

	var roomAdmins []model.RoomAdmin
	roomAdmins = append(roomAdmins, model.RoomAdmin{
		Id: -1,
	})

	item := &memcache.Item{Key: key, Object: roomAdmins, Expiration: d.RoomAdminExpire, Flags: memcache.FlagJSON | memcache.FlagGzip}
	if err = conn.Set(item); err != nil {
		prom.BusinessErrCount.Incr("mc:AddCacheRoomAdminAnchor")
		log.Errorv(c, log.KV("AddCacheRoomAdminAnchor", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

//
//// GetByUserMC .
//func (d *Dao) GetByUserMC(c context.Context, uid int64) (rst []*model.RoomAdmin, hit int64, err error) {
//	hit = 0
//	key := d.getUserKey(uid)
//	conn := d.mc.Get(c)
//	defer conn.Close()
//
//	reply, err := conn.Get(key)
//	//spew.Dump("GetByUserMC1", reply, err)
//
//	if err != nil {
//		if err == memcache.ErrNotFound {
//			err = nil
//			return
//		}
//		log.Error("GetByUserMC get (%v) error(%v)", key, err)
//		return
//	}
//
//	hit = 1
//
//	if reply.Object == nil {
//		return nil, hit, err
//	}
//
//	if err = conn.Scan(reply, rst); err != nil {
//		log.Error("GetByUserMC Scan (%+v) error(%v)", reply, err)
//	}
//
//	//spew.Dump("GetByUserMC2", rst, err)
//	return
//}
//
//// SetByUserMc .
//func (d *Dao) SetByUserMc(c context.Context, uid int64, rst []*model.RoomAdmin) (err error) {
//	key := d.getUserKey(uid)
//	conn := d.mc.Get(c)
//	defer conn.Close()
//
//	err = conn.Set(&memcache.Item{
//		Key:        key,
//		Object:     rst,
//		Flags:      memcache.FlagJSON,
//		Expiration: mcExpire,
//	})
//
//	//spew.Dump("SetByUserMc", err)
//	if err != nil {
//		log.Error("SetByUserMc set(%v) value (%+v) error (%v)", key, rst, err)
//		return
//	}
//	return
//}
