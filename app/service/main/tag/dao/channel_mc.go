package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixChannelGroup = "cg_%d" // key: cg_tid value: []*model.ChannelGroup
)

func keyChannelGroup(tid int64) string {
	return fmt.Sprintf(_prefixChannelGroup, tid)
}

// AddChannelGroupCache add channel group cache.
func (d *Dao) AddChannelGroupCache(c context.Context, tid int64, cgs []*model.ChannelGroup) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{
		Key:        keyChannelGroup(tid),
		Object:     cgs,
		Expiration: d.channelGroupExpire,
		Flags:      memcache.FlagGzip | memcache.FlagJSON,
	}
	if err = conn.Set(item); err != nil {
		log.Error("d.dao.AddChannelGroupCache(%d,%v) error(%v)", tid, item, err)
	}
	return
}

// ChannelGroupCache get channel group cache.
func (d *Dao) ChannelGroupCache(c context.Context, tid int64) (res []*model.ChannelGroup, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	res = make([]*model.ChannelGroup, 0, model.ChannelMaxGroups)
	var (
		item *memcache.Item
		key  = keyChannelGroup(tid)
	)
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			return nil, nil
		}
		log.Error("d.dao.ChannelGroupCache(%d) error(%v)", tid, err)
		return
	}
	if err = conn.Scan(item, &res); err != nil {
		log.Error("d.dao.ChannelGroupCache(%d) conn.Scan(%v) error(%v)", tid, item, err)
	}
	return
}
