package dao

import (
	"context"
	"time"

	"go-common/app/job/main/vip/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Dao struct info of Dao.
type Dao struct {
	// mysql
	db    *sql.DB
	oldDb *sql.DB
	// http
	client *bm.Client
	// conf
	c *conf.Config
	// memcache
	mc       *memcache.Pool
	mcExpire int32
	//redis pool
	redis        *redis.Pool
	redisExpire  int32
	errProm      *prom.Prom
	frozenExpire int32
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// mc
		mc:       memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
		// db
		db:    sql.NewMySQL(c.NewMysql),
		oldDb: sql.NewMySQL(c.OldMysql),
		// http client
		client:       bm.NewClient(c.HTTPClient),
		errProm:      prom.BusinessErrCount,
		frozenExpire: int32(time.Duration(c.Property.FrozenExpire) / time.Second),
	}
	return
}

// Ping ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
}

//StartTx start tx
func (d *Dao) StartTx(c context.Context) (tx *sql.Tx, err error) {
	if d.db != nil {
		tx, err = d.db.Begin(c)
	}
	return
}
