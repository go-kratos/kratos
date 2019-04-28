package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixResource = "r_%d_%d_%d" // key:r_oid_tid_type value:resource
	_prefixResTag   = "rt_%d_%d"
)

func keyResource(oid, tid int64, typ int32) string {
	return fmt.Sprintf(_prefixResource, oid, tid, typ)
}

func keyResTag(oid int64, typ int32) string {
	return fmt.Sprintf(_prefixResTag, oid, typ)
}

// AddResTagCache add resource tag cache.
func (d *Dao) AddResTagCache(c context.Context, oid int64, tp int32, rs []*model.Resource) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{
		Key:        keyResTag(oid, tp),
		Object:     rs,
		Flags:      memcache.FlagJSON,
		Expiration: d.resExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("d.dao.AddResTagCache(%d,%d,%v) error(%v)", oid, tp, rs, err)
	}
	conn.Close()
	return
}

// AddResTagMapCaches add resource tag map caches.
func (d *Dao) AddResTagMapCaches(c context.Context, tp int32, rsMap map[int64][]*model.Resource) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for oid, rs := range rsMap {
		item := &memcache.Item{
			Key:        keyResTag(oid, tp),
			Object:     rs,
			Flags:      memcache.FlagJSON,
			Expiration: d.resExpire,
		}
		if err = conn.Set(item); err != nil {
			log.Error("d.dao.AddResTagMapCache(%d,%d,%v) error(%v)", oid, tp, item, err)
			return
		}
	}
	return
}

// ResTagCache res tag cache.
func (d *Dao) ResTagCache(c context.Context, oid int64, tp int32) (res []*model.Resource, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	res = make([]*model.Resource, 0)
	var (
		item *memcache.Item
		key  = keyResTag(oid, tp)
	)
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			return nil, nil
		}
		log.Error("d.dao.ResTagCache(%d,%d) error(%v)", oid, tp, err)
		return
	}
	if err = conn.Scan(item, &res); err != nil {
		log.Error("d.dao.ResTagCache(%d,%d) conn.Scan(%v) error(%v)", oid, tp, item, err)
	}
	return
}

// ResTagMapCache res tag cache.
func (d *Dao) ResTagMapCache(c context.Context, oid int64, tp int32) (res map[int64]*model.Resource, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	rs := make([]*model.Resource, 0)
	var (
		item *memcache.Item
		key  = keyResTag(oid, tp)
	)
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			return nil, nil
		}
		log.Error("d.dao.ResTagCache(%d,%d) error(%v)", oid, tp, err)
		return
	}
	if err = conn.Scan(item, &rs); err != nil {
		log.Error("d.dao.ResTagCache(%d,%d) conn.Scan(%v) error(%v)", oid, tp, item, err)
		return
	}
	res = make(map[int64]*model.Resource, len(rs))
	for _, r := range rs {
		res[r.Tid] = r
	}
	return
}

// ResTagMapCaches res tag map caches.
func (d *Dao) ResTagMapCaches(c context.Context, oids []int64, tp int32) (res map[int64][]*model.Resource, missed []int64, err error) {
	var (
		keys    = make([]string, 0, len(oids))
		keysMap = make(map[string]int64, len(oids))
		items   map[string]*memcache.Item
	)
	for _, oid := range oids {
		key := keyResTag(oid, tp)
		keys = append(keys, key)
		keysMap[key] = oid
	}
	res = make(map[int64][]*model.Resource, len(oids))
	conn := d.mc.Get(c)
	defer conn.Close()
	if items, err = conn.GetMulti(keys); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			missed = oids
		} else {
			log.Error("d.dao.ResTagMapCache(%v) error(%v)", keys, err)
		}
		return
	}
	for _, item := range items {
		oid, ok := keysMap[item.Key]
		if !ok || oid <= 0 {
			continue
		}
		rs := make([]*model.Resource, 0)
		if err = conn.Scan(item, &rs); err != nil {
			log.Error("d.dao.ResTagMapCache(%v,%d) Scan(%v) error(%v)", oids, tp, item, err)
			return
		}
		delete(keysMap, item.Key)
		res[oid] = rs
	}
	for _, tid := range keysMap {
		missed = append(missed, tid)
	}
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

// AddResourceCache add a resource cache.
func (d *Dao) AddResourceCache(c context.Context, r *model.Resource) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &memcache.Item{
		Key:        keyResource(r.Oid, r.Tid, r.Type),
		Object:     r,
		Flags:      memcache.FlagProtobuf,
		Expiration: d.resExpire,
	}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%v) error(%v)", item, err)
	}
	return
}

// AddResourcesCache add resources cache.
func (d *Dao) AddResourcesCache(c context.Context, rs []*model.Resource) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, r := range rs {
		item := &memcache.Item{
			Key:        keyResource(r.Oid, r.Tid, r.Type),
			Object:     r,
			Flags:      memcache.FlagProtobuf,
			Expiration: d.resExpire,
		}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s) error(%v)", item.Key, err)
			return
		}
	}
	return
}

// AddResourceMapCache .
func (d *Dao) AddResourceMapCache(c context.Context, rsMap map[int64]*model.Resource) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, r := range rsMap {
		item := &memcache.Item{
			Key:        keyResource(r.Oid, r.Tid, r.Type),
			Object:     r,
			Flags:      memcache.FlagProtobuf,
			Expiration: d.resExpire,
		}
		if err = conn.Set(item); err != nil {
			log.Error("conn.Set(%s) error(%v)", item.Key, err)
			return
		}
	}
	return
}

// AddResourceMapCaches .
func (d *Dao) AddResourceMapCaches(c context.Context, rsMaps map[int64][]*model.Resource) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	for _, rsMap := range rsMaps {
		for _, r := range rsMap {
			item := &memcache.Item{
				Key:        keyResource(r.Oid, r.Tid, r.Type),
				Object:     r,
				Flags:      memcache.FlagProtobuf,
				Expiration: d.resExpire,
			}
			if err = conn.Set(item); err != nil {
				log.Error("conn.Set(%s) error(%v)", item.Key, err)
				return
			}
		}
	}
	return
}

// ResourceMapCache return resource map cache by tids.
func (d *Dao) ResourceMapCache(c context.Context, oid int64, typ int32, tids []int64) (res map[int64]*model.Resource, missed []int64, err error) {
	var (
		keys    = make([]string, 0, len(tids))
		keysMap = make(map[string]int64, len(tids))
		items   map[string]*memcache.Item
	)
	for _, tid := range tids {
		key := keyResource(oid, tid, typ)
		keys = append(keys, key)
		keysMap[key] = tid
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	if items, err = conn.GetMulti(keys); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.GetMulti(%v) error(%v)", keys, err)
		}
		return
	}
	res = make(map[int64]*model.Resource, len(items))
	for _, item := range items {
		m := &model.Resource{}
		if err = conn.Scan(item, m); err != nil {
			log.Error("conn.Scan(%s) error(%v)", item.Value, err)
			return
		}
		delete(keysMap, item.Key)
		res[m.Tid] = m
	}
	for _, tid := range keysMap {
		missed = append(missed, tid)
	}
	return
}

// ResourceCache return resource  cache by tid.
func (d *Dao) ResourceCache(c context.Context, oid int64, typ int32, tid int64) (res *model.Resource, err error) {
	var (
		key  = keyResource(oid, tid, typ)
		item *memcache.Item
	)
	conn := d.mc.Get(c)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Get(%v) error(%v)", key, err)
		}
		return
	}
	res = &model.Resource{}
	if err = conn.Scan(item, res); err != nil {
		log.Error("conn.Scan(%s) error(%v)", item.Value, err)
	}
	return
}
