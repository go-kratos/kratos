package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/ugcpay/conf"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

func elecOrderIDKey(orderID string) string {
	return fmt.Sprintf("ur_eoi_%s", orderID)
}

// AddCacheOrderID Set data to mc
func (d *Dao) AddCacheOrderID(c context.Context, orderID string) (ok bool, err error) {
	ok = true
	conn := d.mc.Get(c)
	defer conn.Close()
	key := elecOrderIDKey(orderID)
	item := &memcache.Item{Key: key, Object: struct{}{}, Expiration: conf.Conf.CacheTTL.ElecOrderIDTTL, Flags: memcache.FlagJSON}
	if err = conn.Add(item); err != nil {
		if err == memcache.ErrNotStored {
			ok = false
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:AddCacheOrderID")
		log.Errorv(c, log.KV("AddCacheOrderID", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// DelCacheOrderID .
func (d *Dao) DelCacheOrderID(c context.Context, orderID string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := elecOrderIDKey(orderID)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:DelCacheOrderID")
		log.Errorv(c, log.KV("DelCacheOrderID", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
