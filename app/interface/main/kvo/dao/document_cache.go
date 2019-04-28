package dao

import (
	"context"
	"encoding/json"
	"fmt"

	mc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_docKeyPrefix = "d"
)

// pingMemcache check memcache health
func (d *Dao) pingMemcache(ctx context.Context) (err error) {
	conn := d.cache.Get(ctx)
	item := mc.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}
	err = conn.Set(&item)
	conn.Close()
	return
}

// DocumentCache memcache get document(hash)
func (d *Dao) DocumentCache(ctx context.Context, checkSum int64) (data json.RawMessage, err error) {
	var result *mc.Item
	conn := d.cache.Get(ctx)
	defer conn.Close()
	if result, err = conn.Get(cacheDocKey(checkSum)); err != nil {
		if err == mc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", cacheDocKey(checkSum), err)
		return
	}
	data = json.RawMessage(result.Value)
	return
}

// DelDocumentCache remove memcache document detail
func (d *Dao) DelDocumentCache(ctx context.Context, checkSum int64) (err error) {
	conn := d.cache.Get(ctx)
	defer conn.Close()
	if err = conn.Delete(cacheDocKey(checkSum)); err == mc.ErrNotFound {
		err = nil
	}
	return
}

// SetDocumentCache add document cache
func (d *Dao) SetDocumentCache(ctx context.Context, checkSum int64, data json.RawMessage) (err error) {
	conn := d.cache.Get(ctx)
	defer conn.Close()
	if err = conn.Set(&mc.Item{
		Key:        cacheDocKey(checkSum),
		Value:      data,
		Expiration: d.mcExpire,
	}); err != nil {
		log.Error("dao.SetDocumentCache(%v,%s) err:%v", cacheDocKey(checkSum), data, err)
		return
	}
	return
}

func cacheDocKey(checkSum int64) string {
	return fmt.Sprintf("%v_%v", _docKeyPrefix, checkSum)
}
