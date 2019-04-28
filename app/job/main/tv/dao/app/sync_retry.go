package app

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/tv/model/common"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

// SetRetry save retry model to memcache
func (d *Dao) SetRetry(c context.Context, retry *common.SyncRetry) (err error) {
	var conn = d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: retry.MCKey(), Object: retry, Flags: memcache.FlagJSON, Expiration: d.mcMediaExpire}); err != nil {
		log.Error("conn.Set error(%v)", err)
		return
	}
	return
}

// GetRetry gets retry times
func (d *Dao) GetRetry(c context.Context, req *common.SyncRetry) (times int, err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
	)
	conn = d.mc.Get(c)
	defer conn.Close()
	if item, err = conn.Get(req.MCKey()); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return // 0
		}
		log.Error("GetRetry Req %s, Err %v", req.MCKey(), err)
		return
	}
	if err = json.Unmarshal(item.Value, &req); err != nil {
		log.Error("GetRetry Req %s, Json Err %v", req.MCKey(), err)
		return
	}
	times = req.Retry
	return
}

// DelRetry deletes the cache from MC
func (d *Dao) DelRetry(c context.Context, req *common.SyncRetry) (err error) {
	var (
		key  = req.MCKey()
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		log.Error("conn.Set error(%v)", err)
	}
	return
}
