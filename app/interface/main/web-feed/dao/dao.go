package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/web-feed/conf"
	"go-common/library/cache/memcache"
	"go-common/library/stat/prom"
)

// Dao .
type Dao struct {
	c            *conf.Config
	mc           *memcache.Pool
	mcFeedExpire int32
}

var (
	errorsCount = prom.BusinessErrCount
	infosCount  = prom.BusinessInfoCount
)

// PromError prom error
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo add prom info
func PromInfo(name string) {
	infosCount.Incr(name)
}

// New add a feed job dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		mc:           memcache.NewPool(c.Memcache.Config),
		mcFeedExpire: int32(time.Duration(c.Memcache.FeedExpire) / time.Second),
	}
	return
}

// Ping checks health of redis and mc.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.pingMC(c)
}

// Close closes connections of redis or mc etc.
func (d *Dao) Close() {
	if d.mc != nil {
		d.mc.Close()
	}
}
