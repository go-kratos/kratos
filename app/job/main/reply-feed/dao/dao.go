package dao

import (
	"context"
	"time"

	"go-common/app/job/main/reply-feed/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao
type Dao struct {
	c        *conf.Config
	mc       *memcache.Pool
	mcExpire int32

	redis                *redis.Pool
	redisReplySetExpire  int
	redisReplyZSetExpire int
	redisRefreshExpire   int
	db                   *xsql.DB
	dbSlave              *xsql.DB

	httpCli *bm.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:                    c,
		mc:                   memcache.NewPool(c.Memcache),
		mcExpire:             int32(time.Duration(c.MemcacheExpire.McStatExpire) / time.Second),
		redis:                redis.NewPool(c.Redis),
		redisReplySetExpire:  int(time.Duration(c.RedisExpire.RedisReplySetExpire) / time.Second),
		redisReplyZSetExpire: int(time.Duration(c.RedisExpire.RedisReplyZSetExpire) / time.Second),
		redisRefreshExpire:   int(time.Duration(c.RedisExpire.RedisRefreshExpire) / time.Second),
		db:                   xsql.NewMySQL(c.MySQL.DB),
		dbSlave:              xsql.NewMySQL(c.MySQL.DBSlave),
		httpCli:              bm.NewClient(c.HTTPClient),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.dbSlave.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	if err := d.PingRedis(c); err != nil {
		return err
	}
	if err := d.PingMc(c); err != nil {
		return err
	}
	return d.db.Ping(c)
}
