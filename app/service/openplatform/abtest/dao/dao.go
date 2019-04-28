package dao

import (
	"context"
	"time"

	"go-common/app/service/openplatform/abtest/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao struct answer history of Dao
type Dao struct {
	c *conf.Config
	// db
	db *sql.DB
	// redis
	redis        *redis.Pool
	expire       int
	verifyExpire int
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		db:           sql.NewMySQL(c.DB.Ab),
		redis:        redis.NewPool(c.Redis.Config),
		expire:       int(time.Duration(c.Redis.Expire) / time.Second),
		verifyExpire: int(time.Duration(c.Redis.VerifyCdTimes) / time.Second),
	}
	return
}

// Close close connections.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Error("PingDb error(%v)", err)
		return
	}
	if err = d.PingRedis(c); err != nil {
		log.Error("PingRedis error(%v)", err)
		return
	}
	return
}
