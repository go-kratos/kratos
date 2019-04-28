package dao

import (
	"context"
	"fmt"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

func (d *Dao) Lock(c context.Context, key string, val string) (err error) {
	if val == "" {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	obj := &struct {
		Value string
	}{
		Value: val,
	}
	item := &memcache.Item{Key: key, Object: obj, Expiration: d.cacheTTL.LockTTL, Flags: memcache.FlagJSON}
	if err = conn.Add(item); err != nil {
		prom.BusinessErrCount.Incr("mc:Lock")
		log.Errorv(c, log.KV("Lock", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// no thread-safe, but it's works.
// bad case of unlocking：
// 1, process-a gets lock
// 2, lock expires
// 3, process-b gets lock
// 4, process-a releases lock
func (d *Dao) Unlock(c context.Context, key string, val string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:Unlock")
		log.Errorv(c, log.KV("Unlock", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	res := &struct {
		Value string
	}{}
	err = conn.Scan(reply, &res)
	if err != nil {
		prom.BusinessErrCount.Incr("mc:Unlock")
		log.Errorv(c, log.KV("Unlock", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	if res.Value != val {
		return nil
	}
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:Unlock")
		log.Errorv(c, log.KV("Unlock", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
