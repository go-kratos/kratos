package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/ugcpay/conf"
	"go-common/app/service/main/ugcpay/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

var _ _mc

// CacheOrderUser get data from mc
func (d *Dao) CacheOrderUser(c context.Context, id string) (res *model.Order, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := orderKey(id)
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:CacheOrderUser")
		log.Errorv(c, log.KV("CacheOrderUser", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	res = &model.Order{}
	err = conn.Scan(reply, res)
	if err != nil {
		prom.BusinessErrCount.Incr("mc:CacheOrderUser")
		log.Errorv(c, log.KV("CacheOrderUser", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// AddCacheOrderUser Set data to mc
func (d *Dao) AddCacheOrderUser(c context.Context, id string, val *model.Order) (err error) {
	if val == nil {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	key := orderKey(id)
	item := &memcache.Item{Key: key, Object: val, Expiration: conf.Conf.CacheTTL.OrderTTL, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		prom.BusinessErrCount.Incr("mc:AddCacheOrderUser")
		log.Errorv(c, log.KV("AddCacheOrderUser", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// DelCacheOrderUser delete data from mc
func (d *Dao) DelCacheOrderUser(c context.Context, id string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := orderKey(id)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:DelCacheOrderUser")
		log.Errorv(c, log.KV("DelCacheOrderUser", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// CacheAsset get data from mc
func (d *Dao) CacheAsset(c context.Context, id int64, otype string, currency string) (res *model.Asset, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := assetKey(id, otype, currency)
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:CacheAsset")
		log.Errorv(c, log.KV("CacheAsset", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	res = &model.Asset{}
	err = conn.Scan(reply, res)
	if err != nil {
		prom.BusinessErrCount.Incr("mc:CacheAsset")
		log.Errorv(c, log.KV("CacheAsset", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// AddCacheAsset Set data to mc
func (d *Dao) AddCacheAsset(c context.Context, id int64, otype string, currency string, value *model.Asset) (err error) {
	if value == nil {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	key := assetKey(id, otype, currency)
	item := &memcache.Item{Key: key, Object: value, Expiration: conf.Conf.CacheTTL.AssetTTL, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		prom.BusinessErrCount.Incr("mc:AddCacheAsset")
		log.Errorv(c, log.KV("AddCacheAsset", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// DelCacheAsset delete data from mc
func (d *Dao) DelCacheAsset(c context.Context, id int64, otype string, currency string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := assetKey(id, otype, currency)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:DelCacheAsset")
		log.Errorv(c, log.KV("DelCacheAsset", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
