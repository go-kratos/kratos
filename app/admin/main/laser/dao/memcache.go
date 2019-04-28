package dao

import (
	"context"
	"go-common/app/admin/main/laser/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"strconv"
)

const (
	_prefix = "taskinfo_"
)

func keyTaskInfo(mid int64) string {
	return _prefix + strconv.FormatInt(mid, 10)
}

// TaskInfoCache get taskInfo cache
func (d *Dao) TaskInfoCache(c context.Context, mid int64) (ti *model.TaskInfo, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	r, err = conn.Get(keyTaskInfo(mid))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn Get2(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &ti); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		ti = nil
	}
	return
}

// AddTaskInfoCache add taskInfo cache
func (d *Dao) AddTaskInfoCache(c context.Context, mid int64, ti *model.TaskInfo) (err error) {
	var (
		key = keyTaskInfo(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: ti, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// RemoveTaskInfoCache  remove taskInfo cache
func (d *Dao) RemoveTaskInfoCache(c context.Context, mid int64) (err error) {
	var (
		key = keyTaskInfo(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		log.Error("memcache.Delete(%v) error(%v)", key, err)
	}
	return
}
