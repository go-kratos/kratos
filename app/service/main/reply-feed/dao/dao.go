package dao

import (
	"context"
	"time"

	"go-common/app/service/main/reply-feed/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c *conf.Config

	redis                *redis.Pool
	redisReplyZSetExpire int
	db                   *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:                    c,
		redis:                redis.NewPool(c.Redis),
		redisReplyZSetExpire: int(time.Duration(c.RedisExpire.RedisReplyZSetExpire) / time.Second),
		db:                   xsql.NewMySQL(c.MySQL),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	if err := d.PingRedis(c); err != nil {
		return err
	}
	return d.db.Ping(c)
}
