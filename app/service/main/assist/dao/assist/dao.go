package assist

import (
	"context"
	"time"

	"go-common/app/service/main/assist/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
)

// Dao is assist dao.
type Dao struct {
	// config
	c *conf.Config
	// db
	db *sql.DB
	// mc
	mc       *memcache.Pool
	mcSubExp int32
	// redis
	redis       *redis.Pool
	redisExpire int32
}

// New init api url
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:        c,
		db:       sql.NewMySQL(c.DB.Assist),
		mc:       memcache.NewPool(c.Memcache.Assist.Config),
		mcSubExp: int32(time.Duration(c.Memcache.Assist.SubmitExpire) / time.Second),
		// redis
		redis:       redis.NewPool(c.Redis.Assist.Config),
		redisExpire: int32(time.Duration(c.Redis.Assist.Expire) / time.Second),
	}
	return
}

// Ping include db, mc, redis
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		return
	}
	if err = d.pingMemcache(c); err != nil {
		return
	}
	return
}

// Close include db, mc, redis
func (d *Dao) Close() (err error) {
	if d.db != nil {
		d.db.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	return
}
