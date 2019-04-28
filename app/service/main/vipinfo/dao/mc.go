package dao

import (
	"context"
	"strconv"

	"go-common/app/service/main/vipinfo/model"
	mc "go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

const (
	_prefixInfo   = "i:"
	_prefixFrozen = "vipfrozen:"
)

func keyInfo(mid int64) string {
	return _prefixInfo + strconv.FormatInt(mid, 10)
}

func keyFrozen(mid int64) string {
	return _prefixFrozen + strconv.FormatInt(mid, 10)
}

// pingMC ping memcache.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	err = conn.Set(&mc.Item{
		Key:   "ping",
		Value: []byte("pong"),
	})
	return
}

// CacheInfo get vip info cache.
func (d *Dao) CacheInfo(c context.Context, mid int64) (v *model.VipUserInfo, err error) {
	key := keyInfo(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao cache info")
		return
	}
	v = new(model.VipUserInfo)
	if err = conn.Scan(r, v); err != nil {
		err = errors.Wrapf(err, "dao cache scan info")
	}
	return
}

// AddCacheInfo add vip info cache.
func (d *Dao) AddCacheInfo(c context.Context, mid int64, v *model.VipUserInfo) (err error) {
	item := &mc.Item{
		Key:        keyInfo(mid),
		Object:     v,
		Expiration: d.mcExpire,
		Flags:      mc.FlagProtobuf,
	}
	conn := d.mc.Get(c)
	err = conn.Set(item)
	conn.Close()
	if err != nil {
		err = errors.Wrapf(err, "dao add cache info")
	}
	return
}

// CacheInfos multi get account info from cache.
func (d *Dao) CacheInfos(c context.Context, mids []int64) (res map[int64]*model.VipUserInfo, err error) {
	keys := make([]string, 0, len(mids))
	keyMidMap := make(map[string]int64, len(mids))
	for _, mid := range mids {
		key := keyInfo(mid)
		if _, ok := keyMidMap[key]; !ok {
			// duplicate mid
			keyMidMap[key] = mid
			keys = append(keys, key)
		}
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	rs, err := conn.GetMulti(keys)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao infos cache")
		return
	}
	res = make(map[int64]*model.VipUserInfo, len(mids))
	for _, r := range rs {
		ai := &model.VipUserInfo{}
		conn.Scan(r, ai)
		res[keyMidMap[r.Key]] = ai
	}
	return
}

// AddCacheInfos set account infos cache.
func (d *Dao) AddCacheInfos(c context.Context, vs map[int64]*model.VipUserInfo) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for mid, i := range vs {
		item := &mc.Item{
			Key:        keyInfo(mid),
			Object:     i,
			Flags:      mc.FlagProtobuf,
			Expiration: d.mcExpire,
		}
		err = conn.Set(item)
		if err != nil {
			err = errors.Wrap(err, "dao add infos cache")
		}
	}
	return
}

// CacheVipFrozens multi get vip frozens from cache.
func (d *Dao) CacheVipFrozens(c context.Context, mids []int64) (res map[int64]int, err error) {
	keys := make([]string, 0, len(mids))
	keyMidMap := make(map[string]int64, len(mids))
	for _, mid := range mids {
		key := keyFrozen(mid)
		if _, ok := keyMidMap[key]; !ok {
			// duplicate mid
			keyMidMap[key] = mid
			keys = append(keys, key)
		}
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	rs, err := conn.GetMulti(keys)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrap(err, "dao frozens cache")
		return
	}
	res = make(map[int64]int, len(mids))
	for _, r := range rs {
		ai := 0
		conn.Scan(r, &ai)
		res[keyMidMap[r.Key]] = ai
	}
	return
}

//CacheVipFrozen get vip frozen flag.
func (d *Dao) CacheVipFrozen(c context.Context, mid int64) (val int, err error) {
	key := keyFrozen(mid)
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao cache frozen")
		return
	}
	if err = conn.Scan(item, &val); err != nil {
		err = errors.Wrapf(err, "dao cache scan frozen")
		return
	}
	return
}

// AddCacheFrozen add cache frozen.
func (d *Dao) AddCacheFrozen(c context.Context, mid int64, vipFrozenFlag int) (err error) {
	item := &mc.Item{
		Key:        keyFrozen(mid),
		Object:     vipFrozenFlag,
		Expiration: d.mcExpire,
		Flags:      mc.FlagJSON,
	}
	conn := d.mc.Get(c)
	err = conn.Set(item)
	conn.Close()
	if err != nil {
		err = errors.Wrapf(err, "dao add cache frozen")
	}
	return
}

// DelInfoCache del vip info cache.
func (d *Dao) DelInfoCache(c context.Context, mid int64) (err error) {
	d.delCache(c, keyInfo(mid))
	return
}

func (d *Dao) delCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == mc.ErrNotFound {
			err = nil
		} else {
			err = errors.Wrapf(err, "del cache")
		}
	}
	return
}
