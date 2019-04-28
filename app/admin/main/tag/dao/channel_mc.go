package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixChannelGroup = "cg_%d" // key: cg_tid value: []*model.ChannelGroup
)

func keyChannelGroup(tid int64) string {
	return fmt.Sprintf(_prefixChannelGroup, tid)
}

// DelChannelGroupCache delete channel group cache.
func (d *Dao) DelChannelGroupCache(c context.Context, tid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(keyChannelGroup(tid)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("d.dao.DelChannelGroupCache(%d) error(%v)", tid, err)
		}
	}
	return
}
