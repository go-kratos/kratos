package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixUUID = "uuid_%d_%s"
	_prefixCD   = "cd_%d_%d"
)

func uuidKey(biz int64, uuid string) string {
	return fmt.Sprintf(_prefixUUID, biz, uuid)
}

func cdKey(app, mid int64) string {
	return fmt.Sprintf(_prefixCD, app, mid)
}

// pingMc ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcUUIDExpire}
	err = conn.Set(&item)
	return
}

// ExistsUUIDCache gets uuid from cache.
func (d *Dao) ExistsUUIDCache(c context.Context, biz int64, uuid string) (exist bool, err error) {
	var (
		conn = d.mc.Get(c)
		key  = uuidKey(biz, uuid)
	)
	defer conn.Close()
	if _, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		PromError("mc:get uuid")
		log.Error("ExistsUUIDCache() conn.Get(%s) error(%v)", key, err)
		return
	}
	exist = true
	return
}

// AddUUIDCache adds uuid cache.
func (d *Dao) AddUUIDCache(c context.Context, biz int64, uuid string) (err error) {
	var (
		conn = d.mc.Get(c)
		key  = uuidKey(biz, uuid)
		item = &memcache.Item{Key: key, Value: []byte{}, Expiration: d.mcUUIDExpire}
	)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		PromError("mc:add uuid")
		log.Error("AddUUIDCache() conn.Set(%+v) error(%v)", item, err)
	}
	return
}

// DelUUIDCache delete uuid cache.
func (d *Dao) DelUUIDCache(c context.Context, biz int64, uuid string) (err error) {
	var (
		conn = d.mc.Get(c)
		key  = uuidKey(biz, uuid)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		PromError("mc:del uuid")
		log.Error("DelUUIDCache(%s) error(%v)", key, err)
	}
	return
}

// ExistsCDCache gets cd from cache.
func (d *Dao) ExistsCDCache(ctx context.Context, app, mid int64) (exist bool, err error) {
	var (
		conn = d.mc.Get(ctx)
		key  = cdKey(app, mid)
	)
	defer conn.Close()
	if _, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		PromError("mc:get cd")
		log.Error("ExistsCDCache() conn.Get(%s) error(%v)", key, err)
		return
	}
	exist = true
	return
}

// AddCDCache adds cd cache.
func (d *Dao) AddCDCache(ctx context.Context, app, mid int64) (err error) {
	var (
		conn = d.mc.Get(ctx)
		key  = cdKey(app, mid)
		item = &memcache.Item{Key: key, Value: []byte{}, Expiration: d.mcCDExpire}
	)
	defer conn.Close()
	if err = conn.Set(item); err != nil {
		PromError("mc:add cd")
		log.Error("AddCDCache() conn.Set(%+v) error(%v)", item, err)
	}
	return
}
