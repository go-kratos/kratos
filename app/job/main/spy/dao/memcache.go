package dao

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_partitionKey = "p_%s_%d"
)

// pingMC ping memcache.
func (d *Dao) pingMC(ctx context.Context) (err error) {
	conn := d.mc.Get(ctx)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 0}); err != nil {
		log.Error("conn.Store(set,ping,1) error(%v)", err)
	}
	return
}

func partitionCacheKey(e string, p int32) string {
	return fmt.Sprintf(_partitionKey, e, p)
}

// OffsetCache .
func (d *Dao) OffsetCache(ctx context.Context, event string, p int32) (offset int64, err error) {
	var (
		key  = partitionCacheKey(event, p)
		conn = d.mc.Get(ctx)
	)
	defer conn.Close()
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(get,%d) error(%v)", p, err)
		return
	}
	if err = conn.Scan(reply, &offset); err != nil {
		log.Error("reply.Scan(%s) error(%v)", string(reply.Value), err)
	}
	return
}

// SetOffsetCache .
func (d *Dao) SetOffsetCache(ctx context.Context, event string, p int32, offset int64) (err error) {
	var (
		key  = partitionCacheKey(event, p)
		conn = d.mc.Get(ctx)
	)
	defer conn.Close()

	if err = conn.Set(&memcache.Item{Key: key, Object: offset, Expiration: 0, Flags: memcache.FlagJSON}); err != nil {
		log.Error("conn.Set(%s,%d) error(%v)", key, offset, err)
	}
	return
}
