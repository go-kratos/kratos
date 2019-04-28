package dao

import (
	"context"

	"go-common/app/job/main/coupon/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao
type Dao struct {
	c      *conf.Config
	db     *xsql.DB
	client *bm.Client
	mc     *memcache.Pool
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:      c,
		db:     xsql.NewMySQL(c.MySQL),
		client: bm.NewClient(c.HTTPClient),
		mc:     memcache.NewPool(c.Memcache),
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
	if err = d.db.Ping(c); err != nil {
		return
	}
	err = d.pingMC(c)
	return
}

// pingMc ping
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 60}
	return conn.Set(&item)
}
