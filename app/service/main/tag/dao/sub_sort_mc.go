package dao

import (
	"context"
	"fmt"
	"go-common/app/service/main/tag/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

var (
	_prefixSubSort = "pss_%d_%d" // prefix sub sort + mid + type
)

func keySubSort(mid int64, typ int) string {
	return fmt.Sprintf(_prefixSubSort, mid, typ)
}

// SubSortCache return tag by mid from cache.
func (d *Dao) SubSortCache(c context.Context, mid int64, typ int) (tids []int64, err error) {
	var (
		key  = keySubSort(mid, typ)
		conn = d.mc.Get(c)
		item *memcache.Item
		b    []byte
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &b); err != nil {
		log.Error("mc.Scan(%s) error(%v)", item.Value, err)
	}
	tids, err = model.SetIndex(b)
	return
}

// AddSubSortCache add a sub tag  to cache.
func (d *Dao) AddSubSortCache(c context.Context, mid int64, typ int, tids []int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{
		Key:        keySubSort(mid, typ),
		Value:      model.Index(tids),
		Flags:      memcache.FlagRAW,
		Expiration: d.tagExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// DelSubSortCache delete sub sort cache.
func (d *Dao) DelSubSortCache(c context.Context, mid int64, typ int) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(keySubSort(mid, typ)); err != nil {
		if err == memcache.ErrNotFound {
			return nil
		}
		log.Error("custome sort conn.Delete(%d%d) error(%v)", mid, typ, err)
	}
	return
}
