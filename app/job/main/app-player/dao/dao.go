package dao

import (
	"context"

	"go-common/app/job/main/app-player/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
)

// Dao is dao.
type Dao struct {
	// mc
	mc *memcache.Pool
	// redis
	redis *redis.Pool
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// mc
		mc: memcache.NewPool(c.Memcache),
		// reids
		redis: redis.NewPool(c.Redis),
	}
	return
}

// PingMc is
func (d *Dao) PingMc(c context.Context) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{Key: "ping", Value: []byte{1}, Flags: memcache.FlagRAW, Expiration: 0}
	err = conn.Set(item)
	conn.Close()
	return
}
