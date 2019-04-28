package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

func userKey(mid int64) (key string) {
	key = fmt.Sprintf("u_%d", mid)
	return
}

func syncBlockTypeID() (key string) {
	return "sync_bt_cur_id"
}

func (d *Dao) SetSyncBlockTypeID(c context.Context, id int64) (err error) {
	var (
		key  = syncBlockTypeID()
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: id, Expiration: 3600 * 24, Flags: memcache.FlagJSON}); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (d *Dao) SyncBlockTypeID(c context.Context) (id int64, err error) {
	var (
		key  = syncBlockTypeID()
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	if err = conn.Scan(item, &id); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (d *Dao) DeleteUserCache(c context.Context, mid int64) (err error) {
	var (
		key  = userKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
