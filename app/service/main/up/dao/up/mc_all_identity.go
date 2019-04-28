package up

import (
	"context"
	"strconv"

	"go-common/app/service/main/up/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefix = "allID_"
)

func keyIdentityAll(mid int64) string {
	return _prefix + strconv.FormatInt(mid, 10)
}

// IdentityAllCache get all up all of identify type cache.
func (d *Dao) IdentityAllCache(c context.Context, mid int64) (st *model.IdentifyAll, err error) {
	var (
		conn = d.mcPool.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	r, err = conn.Get(keyIdentityAll(mid))
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get2(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &st); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", r.Value, err)
		st = nil
	}
	return
}

// AddIdentityAllCache add all of up identity type cache.
func (d *Dao) AddIdentityAllCache(c context.Context, mid int64, st *model.IdentifyAll) (err error) {
	var (
		key = keyIdentityAll(mid)
	)
	conn := d.mcPool.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: st, Flags: memcache.FlagJSON, Expiration: d.upExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}
