package dao

import (
	"context"

	"go-common/app/job/live/wallet/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

type Dao struct {
	c  *conf.Config
	db *xsql.DB
	mc *memcache.Pool
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:  c,
		db: xsql.NewMySQL(c.DB.Wallet),
		mc: memcache.NewPool(c.Memcache.Wallet),
	}
	return
}

// Close close the resource.
func (dao *Dao) Close() {
	dao.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Error("PingDb error(%v)", err)
		return
	}
	return d.pingMC(c)
}

// pingMc ping
func (d *Dao) pingMC(c context.Context) (err error) {
	item := &memcache.Item{
		Key:   "ping",
		Value: []byte{1},
	}
	conn := d.mc.Get(c)
	err = conn.Set(item)
	conn.Close()
	if err != nil {
		log.Error("PingMemcache conn.Set(%v) error(%v)", item, err)
	}
	return
}
