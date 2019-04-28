package dao

import (
	"context"
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixDmDailyLimit = "dm_daily_limit_"
)

func keyDmDailyLimit(mid int64) string {
	return _prefixDmDailyLimit + strconv.FormatInt(mid, 10)
}

// SetDmDailyLimitCache .
func (d *Dao) SetDmDailyLimitCache(c context.Context, mid int64, limiter *model.DailyLimiter) (err error) {
	var (
		conn = d.dmSegMC.Get(c)
		key  = keyDmDailyLimit(mid)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     limiter,
		Flags:      memcache.FlagJSON,
		Expiration: d.dmLimiterMCExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// GetDmDailyLimitCache .
func (d *Dao) GetDmDailyLimitCache(c context.Context, mid int64) (limiter *model.DailyLimiter, err error) {
	var (
		conn = d.dmSegMC.Get(c)
		key  = keyDmDailyLimit(mid)
		rp   *memcache.Item
	)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			limiter = nil
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	limiter = &model.DailyLimiter{}
	if err = conn.Scan(rp, &limiter); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}
