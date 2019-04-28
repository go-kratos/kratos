package thirdp

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

func mangoMCKey(isPGC bool, sid int64) string {
	if isPGC {
		return fmt.Sprintf("%s_%d", "mango_pgc", sid)
	}
	return fmt.Sprintf("%s_%d", "mango_ugc", sid)
}

// SetSnCnt save season/archive count
func (d *Dao) SetSnCnt(c context.Context, isPGC bool, sid int64, cnt int) (err error) {
	var (
		key  = mangoMCKey(isPGC, sid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: cnt, Flags: memcache.FlagJSON, Expiration: d.cntExpire}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// GetSnCnt get season/archive count cache.
func (d *Dao) GetSnCnt(c context.Context, isPGC bool, sid int64) (cnt int, err error) {
	var (
		key  = mangoMCKey(isPGC, sid)
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	if item, err = conn.Get(key); err != nil {
		if err != memcache.ErrNotFound {
			log.Error("conn.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(item, &cnt); err != nil {
		log.Error("conn.Get(%s) error(%v)", key, err)
	}
	return
}

// LoadSnCnt loads season's ep cnt, or archive's video cnt
func (d *Dao) LoadSnCnt(ctx context.Context, isPGC bool, sid int64) (cnt int, err error) {
	if cnt, err = d.GetSnCnt(ctx, isPGC, sid); err == nil { // not found or MC error
		return
	}
	log.Warn("LoadSnCnt IsPGC %v, Get Sid [%d] from MC Err (%v)", isPGC, sid, err) // cache set/get error
	if cnt, err = d.MangoSnCnt(ctx, isPGC, sid); err != nil {
		return
	}
	d.addCache(func() {
		d.SetSnCnt(ctx, isPGC, sid, cnt)
	})
	return
}
