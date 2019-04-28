package assist

import (
	"context"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"strconv"
)

func assistKey(mid int64) string {
	return "assist_relation_mid_" + strconv.FormatInt(mid, 10)
}

// SetCacheAss fn
func (d *Dao) SetCacheAss(c context.Context, mid int64, assistMids []int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := assistKey(mid)
	if err = conn.Set(&memcache.Item{Key: key, Object: assistMids, Flags: memcache.FlagJSON, Expiration: 0}); err != nil {
		log.Error("conn.Store error(%v) | key(%s) mid(%d) assistMids(%s)", err, key, mid, assistMids)
	}
	conn.Close()
	return
}

// GetCacheAss GetCacheAss
func (d *Dao) GetCacheAss(c context.Context, mid int64) (assistMids []int64, err error) {
	var (
		key  = assistKey(mid)
		rp   *memcache.Item
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	rp, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			log.Info("memcache.Get(%s) ErrNotFound(%v)", key, err)
			err = nil
		} else {
			log.Error("conn.Get error(%v) | key(%s)", err, key)
		}
		return
	}
	if err = conn.Scan(rp, &assistMids); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", rp.Value, err)
		assistMids = nil
	}
	return
}

// DelCacheAss fn
func (d *Dao) DelCacheAss(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := assistKey(mid)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			log.Warn("DelCacheAss memcache.ErrNotFound key(%s)|(%+v)", key, err)
			err = nil
		} else {
			log.Error("DelCacheAss key(%s)|err(%+v)", key, err)
		}
	}
	return
}

func (d *Dao) pingMemcache(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("mc.ping.Store error(%v)", err)
		return
	}
	return
}
