package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/captcha/conf"
	"go-common/library/cache/memcache"
)

// Dao captcha service Dao.
type Dao struct {
	conf     *conf.Config
	memcache *memcache.Pool
	mcExpire int32
}

// New new a captcha dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:     c,
		memcache: memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
	}
	return d
}

// Ping captcha service health check, connection is ok.
func (d *Dao) Ping(c context.Context) (err error) {
	// return d.pingRedis(c)
	return d.pingMemcache(c)
}

// pingMemcache check Memcache health.
func (d *Dao) pingMemcache(c context.Context) (err error) {
	conn := d.memcache.Get(c)
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}
	err = conn.Set(&item)
	conn.Close()
	return
}

// Close close captcha all connection.
func (d *Dao) Close() {
	if d.memcache != nil {
		d.memcache.Close()
	}
}
