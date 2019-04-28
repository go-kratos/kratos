package dao

import (
	"context"
	"time"

	"go-common/app/job/main/spy/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao event dao def.
type Dao struct {
	c             *conf.Config
	db            *sql.DB
	mc            *memcache.Pool
	hbase         *hbase.Client
	redis         *redis.Pool
	httpClient    *bm.Client
	expire        int
	msgUUIDExpire int
}

// New create instance of dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:             c,
		db:            sql.NewMySQL(c.DB),
		mc:            memcache.NewPool(c.Memcache),
		redis:         redis.NewPool(c.Redis.Config),
		httpClient:    bm.NewClient(c.HTTPClient),
		expire:        int(time.Duration(c.Redis.Expire) / time.Second),
		msgUUIDExpire: int(time.Duration(c.Redis.MsgUUIDExpire) / time.Second),
	}
	if c.HBase != nil {
		d.hbase = hbase.NewClient(c.HBase.Config)
	}
	return
}

// Ping check db health.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		return
	}
	if err = d.PingRedis(c); err != nil {
		return
	}
	if d.hbase != nil {
		if err = d.hbase.Ping(c); err != nil {
			return
		}
	}
	return d.db.Ping(c)
}

// Close close all db connections.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
}
