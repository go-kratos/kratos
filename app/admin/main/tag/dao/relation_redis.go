package dao

import (
	"context"
	"fmt"

	"go-common/library/log"
)

const (
	_oidKey = "ot_%d_%d" // sortedset  key:ot_oid_type value:tid score:ctime
	_tidKey = "tt_%d_%d" // sortedset  key:tt_tid_type value:oid score:ctime
)

func tagResKey(tid int64, typ int32) string {
	return fmt.Sprintf(_tidKey, tid, typ)
}

func oidKey(oid int64, typ int32) string {
	return fmt.Sprintf(_oidKey, oid, typ)
}

// DelRelationCache DelRelationCache.
func (d *Dao) DelRelationCache(c context.Context, oid, tid int64, typ int32) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", oidKey(oid, typ)); err != nil {
		log.Error("DEL relation(%d,%d), error(%v)", oid, typ, err)
	}
	// ZREM tid_type oid will be dangerous.
	if _, err = conn.Do("ZREM", tagResKey(tid, typ), oid); err != nil {
		log.Error("ZREM relation(%d,%d,%d), error(%v)", tid, typ, oid, err)
	}
	return
}
