package fav

import (
	"context"
	"time"

	"go-common/app/job/main/favorite/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao favorite dao.
type Dao struct {
	db          *sql.DB
	redis       *redis.Pool
	mc          *memcache.Pool
	redisExpire int
	mcExpire    int32
}

// New new a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.DB.Fav),
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int(time.Duration(c.Redis.Expire) / time.Second),
		// memcache
		mc:       memcache.NewPool(c.Memcache.Config),
		mcExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
	}
	return
}

// Close close all connection.
func (d *Dao) Close() (err error) {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
	return
}

// Ping ping all resource.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		log.Error("d.pingRedis error(%v)", err)
		return
	}
	if err = d.pingMC(c); err != nil {
		log.Error("d.pingMC error(%v)", err)
		return
	}
	if err = d.pingMySQL(c); err != nil {
		log.Error("d.pingMySQL error(%v)", err)
	}
	return
}

// BeginTran crate a *sql.Tx for database transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}
