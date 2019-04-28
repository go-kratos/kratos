package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"

	"github.com/pkg/errors"
)

func orderKey(id string) string {
	return fmt.Sprintf("up_o_%s", id)
}

func assetKey(oid int64, otype string, currency string) string {
	return fmt.Sprintf("up_a_%d_%s_%s", oid, otype, currency)
}

func taskKey(task string) string {
	return fmt.Sprintf("up_t_%s", task)
}

// DelCacheOrderUser delete data from mc
func (d *Dao) DelCacheOrderUser(c context.Context, orderID string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := orderKey(orderID)
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

// DelCacheAsset delete data from mc
func (d *Dao) DelCacheAsset(c context.Context, oid int64, otype string, currency string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := assetKey(oid, otype, currency)
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

// AddCacheTask .
func (d *Dao) AddCacheTask(c context.Context, task string, ttl int32) (ok bool, err error) {
	var (
		conn = d.mc.Get(c)
		item = &memcache.Item{
			Key:        taskKey(task),
			Value:      []byte{0},
			Expiration: ttl,
		}
	)
	defer conn.Close()

	if err = conn.Add(item); err != nil {
		if err == memcache.ErrNotStored {
			err = nil
			ok = false
			return
		}
		err = errors.WithStack(err)
		return
	}
	ok = true
	return
}

// DelCacheTask .
func (d *Dao) DelCacheTask(c context.Context, task string) (err error) {
	var (
		conn = d.mc.Get(c)
		key  = taskKey(task)
	)
	defer conn.Close()

	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}

// ugcpay-rank

func elecUserSetting(mid int64) string {
	return fmt.Sprintf("eus_%d", mid)
}

// DelCacheUserSetting .
func (d *Dao) DelCacheUserSetting(c context.Context, mid int64) (err error) {
	var (
		conn = d.mcRank.Get(c)
		key  = elecUserSetting(mid)
	)
	defer conn.Close()

	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
	}
	return
}
