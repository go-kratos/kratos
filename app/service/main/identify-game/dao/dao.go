package dao

import (
	"context"

	"go-common/app/service/main/identify-game/conf"
	"go-common/library/cache/memcache"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

var (
	errorsCount = prom.BusinessErrCount
	cachedCount = prom.CacheHit
	missedCount = prom.CacheMiss
)

// PromError prom error
func PromError(name string) {
	errorsCount.Incr(name)
}

// Dao struct info of Dao
type Dao struct {
	c *conf.Config

	mc     *memcache.Pool
	client *httpx.Client
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		mc:     memcache.NewPool(c.Memcache),
		client: httpx.NewClient(c.HTTPClient),
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() error {
	return d.mc.Close()
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		PromError("mc:Ping")
	}
	return
}

// pingMc ping memcache
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 100}
	err = conn.Set(&item)
	return
}
