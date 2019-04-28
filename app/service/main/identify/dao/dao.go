package dao

import (
	"context"

	"go-common/app/service/main/identify/conf"
	"go-common/library/cache/memcache"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

const (
	_tokenURI  = "/intranet/auth/tokenInfo"
	_cookieURI = "/intranet/auth/cookieInfo"
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

	tokenURI  string
	cookieURI string
	mc        *memcache.Pool
	mcLogin   *memcache.Pool
	client    *bm.Client
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:         c,
		tokenURI:  c.Identify.AuthHost + _tokenURI,
		cookieURI: c.Identify.AuthHost + _cookieURI,
		mc:        memcache.NewPool(c.Memcache),
		mcLogin:   memcache.NewPool(c.MemcacheLoginLog),
		client:    bm.NewClient(c.HTTPClient),
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	d.mc.Close()
	d.mcLogin.Close()
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
