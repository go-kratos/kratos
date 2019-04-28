package dao

import (
	"context"
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixUpFilter     = "filter_up_"
	_prefixUserFilter   = "filter_user_"
	_prefixGlobalFilter = "filter_global"
)

func keyUserFilter(mid int64) string {
	return _prefixUserFilter + strconv.FormatInt(mid, 10)
}

func keyUpFilter(mid int64) string {
	return _prefixUpFilter + strconv.FormatInt(mid, 10)
}

func keyGlobalFilter() string {
	return _prefixGlobalFilter
}

// AddUserFilterCache set user filters into cache.
func (d *Dao) AddUserFilterCache(c context.Context, mid int64, data []*model.UserFilter) (err error) {
	bs, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal() error(%v)", err)
		return
	}
	key := keyUserFilter(mid)
	conn := d.filterMC.Get(c)
	item := &memcache.Item{
		Key:        key,
		Value:      bs,
		Expiration: d.filterMCExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("memcache.Set(%v) error(%v)", item, err)
	}
	conn.Close()
	return
}

// DelUserFilterCache delete user filters from cache.
func (d *Dao) DelUserFilterCache(c context.Context, mid int64) (err error) {
	var (
		key  = keyUserFilter(mid)
		conn = d.filterMC.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// UserFilterCache get user filters from cache.
func (d *Dao) UserFilterCache(c context.Context, mid int64) (data []*model.UserFilter, err error) {
	var (
		key  = keyUserFilter(mid)
		conn = d.filterMC.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			data = nil
			PromCacheMiss("user_filter", 1)
		} else {
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("user_filter", 1)
	if e := json.Unmarshal(item.Value, &data); e != nil {
		log.Error("json.Unmarshal(%s) error(%v)", item.Value, e)
	}
	return
}

// AddUpFilterCache add upper filter cache.
func (d *Dao) AddUpFilterCache(c context.Context, mid int64, data []*model.UpFilter) (err error) {
	var (
		conn = d.filterMC.Get(c)
		key  = keyUpFilter(mid)
	)
	defer conn.Close()
	item := &memcache.Item{
		Key:        key,
		Object:     data,
		Flags:      memcache.FlagJSON | memcache.FlagGzip,
		Expiration: d.filterMCExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// DelUpFilterCache delete up filters from cache.
func (d *Dao) DelUpFilterCache(c context.Context, mid int64) (err error) {
	key := keyUpFilter(mid)
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

// UpFilterCache get user filter from memcache.
func (d *Dao) UpFilterCache(c context.Context, mid int64) (data []*model.UpFilter, err error) {
	var (
		conn = d.filterMC.Get(c)
		key  = keyUpFilter(mid)
		rp   *memcache.Item
	)
	defer conn.Close()
	data = make([]*model.UpFilter, 0)
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			PromCacheMiss("upper_filter", 1)
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("upper_filter", 1)
	if err = conn.Scan(rp, &data); err != nil {
		log.Error("mc.Scan() error(%v)", err)
	}
	return
}

// AddGlobalFilterCache set global rules into cache.
func (d *Dao) AddGlobalFilterCache(c context.Context, data []*model.GlobalFilter) (err error) {
	var (
		value []byte
		key   = keyGlobalFilter()
		conn  = d.filterMC.Get(c)
	)
	defer conn.Close()
	if value, err = json.Marshal(data); err != nil {
		log.Error("json.Marshal(%v) error(%v)", data, err)
		return
	}
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: d.filterMCExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("memcache.Set(%v) error(%v)", item, err)
	}
	return
}

// DelGlobalFilterCache delete global rules from cache.
func (d *Dao) DelGlobalFilterCache(c context.Context) (err error) {
	var (
		key  = keyGlobalFilter()
		conn = d.filterMC.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("memcache.Delete(%s) error(%v)", key, err)
		}
	}
	return
}

// GlobalFilterCache get up filters from cache.
func (d *Dao) GlobalFilterCache(c context.Context) (data []*model.GlobalFilter, err error) {
	var (
		key  = keyGlobalFilter()
		conn = d.filterMC.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			data = nil
			PromCacheMiss("global_filter", 1)
		} else {
			log.Error("memcache.Get(%s) error(%v)", key, err)
		}
		return
	}
	PromCacheHit("global_filter", 1)
	if e := json.Unmarshal(item.Value, &data); e != nil {
		log.Error("json.Unmarshal(%s) error(%v)", item.Value, e)
	}
	return
}
