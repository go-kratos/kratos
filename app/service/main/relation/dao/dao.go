package dao

import (
	"context"
	"time"

	"github.com/bluele/gcache"

	"go-common/app/service/main/relation/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

const (
	_passportDetailURL = "/intranet/acc/detail"
)

// Dao struct info of Dao.
type Dao struct {
	// mysql
	db *sql.DB
	// memcache
	mc             *memcache.Pool
	followerExpire int32
	mcExpire       int32
	// redis
	redis       *redis.Pool
	redisExpire int32
	// conf
	c *conf.Config
	// prompt
	period int64
	bcount int64
	ucount int64
	// followers unread duration
	UnreadDuration int64
	// apis
	detailURI string
	// client
	client *bm.Client
	// statCache
	statStore gcache.Cache
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// db
		db: sql.NewMySQL(c.Mysql),
		// memcache
		mc:             memcache.NewPool(c.Memcache.Config),
		mcExpire:       int32(time.Duration(c.Memcache.Expire) / time.Second),
		followerExpire: int32(time.Duration(c.Memcache.FollowerExpire) / time.Second),
		// redis
		redis:       redis.NewPool(c.Redis.Config),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
		// prompt
		period: int64(time.Duration(c.Relation.Period) / time.Second),
		bcount: c.Relation.Bcount,
		ucount: c.Relation.Ucount,
		// followers unread
		UnreadDuration: int64(time.Duration(c.Relation.FollowersUnread) / time.Second),
		// passport api
		detailURI: c.Host.Passport + _passportDetailURL,
		client:    bm.NewClient(c.HTTPClient),
		statStore: gcache.New(c.StatCache.Size).LFU().Build(),
	}
	return
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMC(c); err != nil {
		return
	}
	return d.pingRedis(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.mc != nil {
		d.mc.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
}
