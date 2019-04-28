package dao

import (
	"context"

	"go-common/app/job/live-userexp/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao struct userexp-service dao
type Dao struct {
	c *conf.Config
	// exp db
	expDb *sql.DB
	// memcache
	expMc       *memcache.Pool
	cacheExpire int32
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		expDb:       sql.NewMySQL(c.DB.Exp),
		expMc:       memcache.NewPool(c.Memcache.Exp),
		cacheExpire: c.LevelExpire,
	}
	return
}

// Ping check service health.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.expDb.Ping(c); err != nil {
		log.Error("PingDb error(%v)", err)
		return
	}
	if err = d.pingMemcache(c); err != nil {
		return
	}
	return
}

// PingMemcache check connection success.
func (d *Dao) pingMemcache(c context.Context) (err error) {
	item := &memcache.Item{
		Key:        "ping",
		Value:      []byte{1},
		Expiration: d.cacheExpire,
	}
	conn := d.expMc.Get(c)
	err = conn.Set(item)
	conn.Close()
	if err != nil {
		log.Error("PingMemcache conn.Set(%v) error(%v)", item, err)
	}
	return
}

// Close close memcache resource.
func (d *Dao) Close() {
	if d.expMc != nil {
		d.expMc.Close()
	}
	if d.expDb != nil {
		d.expDb.Close()
	}
}
