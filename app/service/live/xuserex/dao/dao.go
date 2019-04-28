package dao

import (
	"context"

	"go-common/app/service/live/xuserex/conf"
	"go-common/library/cache/memcache"
)

// Dao dao
type Dao struct {
	c  *conf.Config
	mc *memcache.Pool
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:  c,
		mc: memcache.NewPool(c.Memcache),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}
