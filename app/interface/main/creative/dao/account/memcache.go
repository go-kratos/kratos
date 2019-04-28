package account

import (
	"context"
	"strconv"

	accmdl "go-common/app/interface/main/creative/model/account"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_upInfoPrefix        = "upinfo_"
	_addMidHalfMinPrefix = "add_midhafmin_"
)

func limitMidHafMin(mid int64) string {
	return _addMidHalfMinPrefix + strconv.FormatInt(mid, 10)
}

func keyUpInfo(mid int64) string {
	return _upInfoPrefix + strconv.FormatInt(mid, 10)
}

// HalfMin fn
func (d *Dao) HalfMin(c context.Context, mid int64) (exist bool, ts uint64, err error) {
	var (
		conn = d.mc.Get(c)
		rp   *memcache.Item
	)
	defer conn.Close()
	key := limitMidHafMin(mid)
	rp, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get error(%v) | key(%s) mid(%d)", err, key, mid)
		}
		return
	}
	if err = conn.Scan(rp, &ts); err != nil {
		log.Error("conn.Scan(%s) error(%v)", rp.Value, err)
		return
	}
	log.Info("HalfMin (%d) | key(%s) ts(%d)", mid, key, ts)
	if ts != 0 {
		exist = true
	}
	return
}

// UpInfoCache get stat cache.
func (d *Dao) UpInfoCache(c context.Context, mid int64) (st *accmdl.UpInfo, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	r, err = conn.Get(keyUpInfo(mid))
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

// AddUpInfoCache add stat cache.
func (d *Dao) AddUpInfoCache(c context.Context, mid int64, st *accmdl.UpInfo) (err error) {
	var (
		key = keyUpInfo(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: st, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}
