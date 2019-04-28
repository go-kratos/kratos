package data

import (
	"context"
	"strconv"

	"go-common/app/interface/main/creative/model/data"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefix       = "s_"
	_upBasePrefix = "sup_"
)

func keyStat(mid int64) string {
	return _prefix + strconv.FormatInt(mid, 10)
}

func keyUpStat(mid int64, date string) string {
	return _upBasePrefix + date + strconv.FormatInt(mid, 10)
}

// statCache get stat cache.
func (d *Dao) statCache(c context.Context, mid int64) (st *data.Stat, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	r, err = conn.Get(keyStat(mid))
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

// addStatCache add stat cache.
func (d *Dao) addStatCache(c context.Context, mid int64, st *data.Stat) (err error) {
	var (
		key = keyStat(mid)
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: st, Flags: memcache.FlagJSON, Expiration: d.mcExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// upBaseStatCache get stat cache.
func (d *Dao) upBaseStatCache(c context.Context, mid int64, dt string) (st *data.UpBaseStat, err error) {
	var (
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	r, err = conn.Get(keyUpStat(mid, dt))
	if err != nil {
		if err == memcache.ErrNotFound {
			log.Error("upBaseStatCache conn.Get2(%d) key not found", mid)
			err = nil
		} else {
			log.Error("upBaseStatCache conn.Get2(%d) error(%v)", mid, err)
		}
		return
	}
	if err = conn.Scan(r, &st); err != nil {
		log.Error("upBaseStatCache json.Unmarshal(%s) error(%v)", r.Value, err)
		st = nil
	}
	return
}

// addUpBaseStatCache add stat cache.
func (d *Dao) addUpBaseStatCache(c context.Context, mid int64, dt string, st *data.UpBaseStat) (err error) {
	key := keyUpStat(mid, dt)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: st, Flags: memcache.FlagJSON, Expiration: d.mcIdxExpire}); err != nil {
		log.Error("memcache.Set(%v) error(%v)", key, err)
	}
	return
}

// DelUpBaseStatCache fn
func (d *Dao) DelUpBaseStatCache(c context.Context, mid int64, dt string) (err error) {
	key := keyUpStat(mid, dt)
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		log.Error("memcache.del(%v) error(%v)", key, err)
	}
	return
}
