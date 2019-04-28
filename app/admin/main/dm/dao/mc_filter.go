package dao

import (
	"context"

	"go-common/library/cache/memcache"
	"go-common/library/log"

	"strconv"
)

const _prefixUpFilter = "filter_up_"

func keyUpFilter(mid, oid int64) string {
	return _prefixUpFilter + strconv.FormatInt(mid, 10) + "_" + strconv.FormatInt(oid, 10)
}

// DelUpFilterCache delete up filters from cache.
func (d *Dao) DelUpFilterCache(c context.Context, mid, oid int64) (err error) {
	key := keyUpFilter(mid, oid)
	conn := d.filterMC.Get(c)
	err = conn.Delete(key)
	conn.Close()
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Delete(%s) error(%v)", key, err)
		}
	}
	return
}
