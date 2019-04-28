package dao

import (
	"context"
	"time"

	"go-common/app/job/main/relation/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao struct info of Dao.
type Dao struct {
	c      *conf.Config
	client *bm.Client
	//path
	clearFollowingPath string
	clearFollowerPath  string
	clearStatPath      string
	// api path
	followersNotify string
	// db
	db *xsql.DB
	// redis
	// redis       *redis.Pool
	// redisExpire int32
	// relation cache
	relRedis  *redis.Pool
	relExpire int32
	mc        *memcache.Pool
	// UnreadDuration int64
}

// New new a Dao and return.
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:                  c,
		client:             bm.NewClient(c.HTTPClient),
		clearFollowingPath: c.ClearPath.Following,
		clearFollowerPath:  c.ClearPath.Follower,
		clearStatPath:      c.ClearPath.Stat,
		followersNotify:    c.ApiPath.FollowersNotify,
		db:                 xsql.NewMySQL(c.Mysql),
		// redis:              redis.NewPool(c.Redis.Config),
		// redisExpire: int32(time.Duration(c.RelRedis.Expire) / time.Second),
		relRedis:  redis.NewPool(c.RelRedis.Config),
		relExpire: int32(time.Duration(c.RelRedis.Expire) / time.Second),
		mc:        memcache.NewPool(c.Memcache.Config),
		// UnreadDuration:     int64(time.Duration(c.Relation.FollowersUnread) / time.Second),
	}
	return dao
}

// Ping ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}
