package dao

import (
	"context"
	"time"

	"go-common/app/service/main/figure/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
)

// Dao figure DAO
type Dao struct {
	c           *conf.Config
	db          *sql.DB
	redis       *redis.Pool
	redisExpire int32
}

// New new a figure DAO
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		db:          sql.NewMySQL(c.Mysql),
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
	}
	return
}

// Ping check service health
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.PingRedis(c); err != nil {
		return
	}
	return d.db.Ping(c)
}

// Close close all dao.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
