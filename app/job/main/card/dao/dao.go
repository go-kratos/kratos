package dao

import (
	"context"

	"go-common/app/job/main/card/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c  *conf.Config
	mc *memcache.Pool
	db *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:  c,
		mc: memcache.NewPool(c.Memcache.Config),
		db: xsql.NewMySQL(c.MySQL),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		return
	}
	return d.db.Ping(c)
}
