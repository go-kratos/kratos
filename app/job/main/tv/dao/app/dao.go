package app

import (
	"time"

	"go-common/app/job/main/tv/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Dao dao.
type Dao struct {
	conf *conf.Config
	// DB
	DB *sql.DB
	// Memcache
	mc            *memcache.Pool
	mcExpire      int32
	mcMediaExpire int32
	// Http client
	client *httpx.Client
	// redis
	redis       *redis.Pool
	redisExpire int32
}

var (
	errorsCount = prom.BusinessErrCount
	infosCount  = prom.BusinessInfoCount
)

// PromError prometheus error count.
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo prometheus info count.
func PromInfo(name string) {
	infosCount.Incr(name)
}

// New create a instance of Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		conf: c,
		// db
		DB: sql.NewMySQL(c.Mysql),
		// mc
		mc:            memcache.NewPool(c.Memcache.Config),
		mcExpire:      int32(time.Duration(c.Memcache.Expire) / time.Second),
		mcMediaExpire: int32(time.Duration(c.Memcache.ExpireMedia) / time.Second),
		client:        httpx.NewClient(conf.Conf.HTTPClient),
		redis:         redis.NewPool(c.Redis.Config),
		redisExpire:   int32(time.Duration(c.Redis.Expire) / time.Second),
	}
	return
}

// Close close the redis and kafka resource.
func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
}

// NumPce calculates number of piece
func NumPce(count int, pagesize int) (numPce int) {
	if count%pagesize == 0 {
		numPce = count / pagesize
		return
	}
	return count/pagesize + 1
}
