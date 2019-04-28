package dao

import (
	"context"
	"runtime"
	"time"

	"go-common/app/interface/main/shorturl/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao struct conf
type Dao struct {
	db       *sql.DB
	memchDB  *memcache.Pool
	mcExpire int32
	cacheCh  chan func()
}

// New new dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:       sql.NewMySQL(c.Mysql),
		memchDB:  memcache.NewPool(c.Memcache.Config),
		cacheCh:  make(chan func(), 1024),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		go d.cacheproc()
	}
	return
}

func (d *Dao) cacheproc() {
	for {
		f, ok := <-d.cacheCh
		if !ok {
			return
		}
		f()
	}
}

// AddCache add cache chan
func (d *Dao) AddCache(f func()) {
	select {
	case d.cacheCh <- f:
	default:
		log.Error("d.cacheCh is full")
	}
}

// Close close db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// Ping ping dao
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Error("d.PingDB error(%v)", err)
	}
	if err = d.PingMC(c); err != nil {
		log.Error("d.PingMC error(%v)", err)
	}
	return
}

// PingMC ping mc is ok.
func (d *Dao) PingMC(c context.Context) (err error) {
	conn := d.memchDB.Get(c)
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}
	if err = conn.Set(&item); err != nil {
		log.Error("conn.Set(%s) error(%v)", item.Key, err)
	}
	conn.Close()
	return
}
