package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/tag/model"
	"go-common/library/log"
)

const (
	_tagResKey = "trt_%d_%d" // SortedSet: key: trt_tid_typ   score: mtime  member: oid
	_oidKey    = "ot_%d_%d"  // sortedset  key:ot_oid_type value:tid score:ctime  // oid cache, will drop.
)

func tagResKey(tid int64, tp int32) string {
	return fmt.Sprintf(_tagResKey, tid, tp)
}

func resOidKey(oid int64, typ int32) string {
	return fmt.Sprintf(_oidKey, oid, typ)
}

// AddTagResCache add tag res cache.
func (d *Dao) AddTagResCache(c context.Context, rt *model.ResTag) (err error) {
	conn := d.redisRank.Get(c)
	defer conn.Close()
	var count int
	for _, tid := range rt.Tids {
		count++
		if err = conn.Send("ZADD", tagResKey(tid, rt.Type), rt.MTime, rt.Oid); err != nil {
			log.Error("d.AddTagResCache(%v)(ZADD, %v, %v, %v) error(%v)", rt, tid, rt.MTime, rt.Oid, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("d.AddTagResCache(%v) Flush() error(%v)", rt, err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("d.AddTagResCache(%v) Receive error(%v)", rt, err)
			return
		}
	}
	return
}

// RemoveTagResCache remove tag res cache.
func (d *Dao) RemoveTagResCache(c context.Context, rt *model.ResTag) (err error) {
	conn := d.redisRank.Get(c)
	defer conn.Close()
	var count int
	for _, tid := range rt.Tids {
		count++
		if err = conn.Send("ZREM", tagResKey(tid, rt.Type), rt.Oid); err != nil {
			log.Error("d.RemoveTagResCache(%v)(ZADD, %v, %v) error(%v)", rt, tid, rt.Oid, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("d.RemoveTagResCache(%v) Flush() error(%v)", rt, err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("d.RemoveTagResCache(%v) Receive error(%v)", rt, err)
			return
		}
	}
	return
}

// DelResOidCache .
func (d *Dao) DelResOidCache(c context.Context, oid int64, tp int32) (err error) {
	conn := d.redisTag.Get(c)
	if _, err = conn.Do("del", resOidKey(oid, tp)); err != nil {
		log.Error("d.DelResOidCache(%d,%d) error(%v)", oid, tp, err)
	}
	conn.Close()
	return
}
