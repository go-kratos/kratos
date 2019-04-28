package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixResource = "r_%d_%d_%d" // key:r_oid_tid_type value:resource
	_prefixResTag   = "rt_%d_%d"   // key:r_oid_type value:[]*resource
)

func keyResource(oid, tid int64, typ int32) string {
	return fmt.Sprintf(_prefixResource, oid, tid, typ)
}
func keyResTag(oid int64, typ int32) string {
	return fmt.Sprintf(_prefixResTag, oid, typ)
}

// DelResMemCache .
func (d *Dao) DelResMemCache(c context.Context, oid, tid int64, tp int32) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(keyResource(oid, tid, tp)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("relation cache conn.Delete(%d,%d,%d) error(%v)", oid, tid, tp, err)
		}
	}
	conn.Close()
	return
}

// DelResTagCache delete res tag cache.
func (d *Dao) DelResTagCache(c context.Context, oid int64, tp int32) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(keyResTag(oid, tp)); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("d.dao.DelResTagCache(%d,%d) error(%v)", oid, tp, err)
		}
	}
	conn.Close()
	return
}
