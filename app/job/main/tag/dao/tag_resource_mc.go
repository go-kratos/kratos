package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixResTag = "rt_%d_%d"
)

func keyResTag(oid int64, typ int32) string {
	return fmt.Sprintf(_prefixResTag, oid, typ)
}

// DelTagResourceCache delete tag resource cache.
func (d *Dao) DelTagResourceCache(c context.Context, oid int64, tp int32) (err error) {
	conn := d.memcache.Get(c)
	if err = conn.Delete(keyResTag(oid, tp)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("d.dao.DelResTagCache(%d,%d) error(%v)", oid, tp, err)
		}
	}
	defer conn.Close()
	return
}
