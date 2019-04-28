package dao

import (
	"context"
	"fmt"

	"go-common/app/service/live/userexp/model"
	mc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_expKey = "level:%d"
)

func key(uid int64) string {
	return fmt.Sprintf(_expKey, uid)
}

// LevelCache 获取等级缓存
func (d *Dao) LevelCache(c context.Context, uid int64) (level *model.Level, err error) {
	key := key(uid)
	conn := d.expMc.Get(c)
	defer conn.Close()

	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		log.Error("[dao.mc_exp|LevelCache] conn.Get(%s) error(%v)", key, err)
		return
	}
	level = &model.Level{}
	if err = conn.Scan(r, level); err != nil {
		log.Error("[dao.mc_exp|LevelCache] conn.Scan(%s) error(%v)", string(r.Value), err)
	}
	return
}

// SetLevelCache 设置等级缓存
func (d *Dao) SetLevelCache(c context.Context, level *model.Level) (err error) {
	key := key(level.Uid)
	conn := d.expMc.Get(c)
	defer conn.Close()

	if conn.Set(&mc.Item{
		Key:        key,
		Object:     level,
		Flags:      mc.FlagProtobuf,
		Expiration: d.cacheExpire,
	}); err != nil {
		log.Error("[dao.mc_exp|SetLevelCache] conn.Set(%s, %v) error(%v)", key, level, err)
	}
	return
}

// DelLevelCache 删除等级缓存
func (d *Dao) DelLevelCache(c context.Context, uid int64) (err error) {
	key := key(uid)
	conn := d.expMc.Get(c)
	defer conn.Close()

	if err = conn.Delete(key); err == mc.ErrNotFound {
		err = nil
	} else if err != nil {
		log.Error("[dao.mc_exp|DelLevelCache] conn.Delete(%s) error(%v)", key, err)
	}
	return
}

// MultiLevelCache 批量获取等级缓存
func (d *Dao) MultiLevelCache(c context.Context, uids []int64) (level []*model.Level, missed []int64, err error) {
	var keys []string
	var um = make(map[int64]struct{}, len(uids))
	for _, uid := range uids {
		keys = append(keys, key(uid))
		um[uid] = struct{}{}
	}
	conn := d.expMc.Get(c)
	defer conn.Close()

	r, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("[dao.mc_exp|MultiLevelCache] conn.GetMulti error(%v)", err)
		return
	}
	// 命中列表
	for _, v := range r {
		ele := &model.Level{}
		if err = conn.Scan(v, ele); err != nil {
			log.Error("[dao.mc_exp|MultiLevelCache] conn.Scan error(%v)", err)
			return
		}
		level = append(level, ele)
		delete(um, ele.Uid)
	}

	// MISS列表
	if len(level) != 0 {
		for uid := range um {
			missed = append(missed, uid)
		}
	} else {
		missed = uids
	}
	return
}
