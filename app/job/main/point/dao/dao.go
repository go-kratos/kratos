package dao

import (
	"context"

	"go-common/app/job/main/point/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao
type Dao struct {
	c      *conf.Config
	mc     *memcache.Pool
	db     *xsql.DB
	client *bm.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:      c,
		mc:     memcache.NewPool(c.Memcache),
		db:     xsql.NewMySQL(c.MySQL),
		client: bm.NewClient(c.HTTPClient),
	}
	return
}

// Close close the resource.
func (dao *Dao) Close() {
	dao.mc.Close()
	dao.db.Close()
}

// Ping dao ping
func (dao *Dao) Ping(c context.Context) (err error) {
	if err = dao.db.Ping(c); err != nil {
		return
	}
	err = dao.pingMC(c)
	return
}

// pingMc ping
func (dao *Dao) pingMC(c context.Context) (err error) {
	conn := dao.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 60}
	return conn.Set(&item)
}
