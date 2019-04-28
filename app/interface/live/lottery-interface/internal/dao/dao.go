package dao

import (
	"context"

	"go-common/app/interface/live/lottery-interface/internal/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c     *conf.Config
	mc    *memcache.Pool
	redis *redis.Pool
	db    *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		mc:    memcache.NewPool(c.Memcache),
		redis: redis.NewPool(c.Redis),
		db:    xsql.NewMySQL(c.MySQL),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return d.db.Ping(ctx)
}
