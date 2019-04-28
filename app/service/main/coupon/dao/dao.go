package dao

import (
	"context"
	"time"

	"go-common/app/service/main/coupon/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Dao dao
type Dao struct {
	c           *conf.Config
	mc          *memcache.Pool
	db          *xsql.DB
	errProm     *prom.Prom
	mcExpire    int32
	prizeExpire int32
	client      *bm.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:           c,
		mc:          memcache.NewPool(c.Memcache.Config),
		db:          xsql.NewMySQL(c.MySQL),
		errProm:     prom.BusinessErrCount,
		mcExpire:    int32(time.Duration(c.Memcache.Expire) / time.Second),
		prizeExpire: int32(time.Duration(c.Memcache.PrizeExpire) / time.Second),
		client:      bm.NewClient(c.HTTPClient),
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
		d.errProm.Incr("ping_db")
		return
	}
	err = d.pingMC(c)
	return
}

// pingMc ping
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}
	return conn.Set(&item)
}
