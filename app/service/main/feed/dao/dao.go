package dao

import (
	"context"
	"time"

	"go-common/app/service/main/feed/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

var (
	// CachedCount .
	CachedCount = prom.CacheHit
	// MissedCount .
	MissedCount = prom.CacheMiss
	infosCount  = prom.BusinessInfoCount
	warnsCount  = prom.BusinessErrCount
)

// PromError stat and log.
func PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}

// PromInfo add prom info
func PromInfo(name string) {
	infosCount.Incr(name)
}

// PromWarn add prom warn
func PromWarn(name string) {
	warnsCount.Incr(name)
}

// Dao struct info of Dao.
type Dao struct {
	// redis
	redis                  *redis.Pool
	redisTTLUpper          int32
	redisExpireUpper       int32
	redisExpireFeed        int32
	redisExpireArchiveFeed int32
	redisExpireBangumiFeed int32
	// memcache
	mc            *memcache.Pool
	mcExpire      int32
	bangumiExpire int32
	// feed Config
	appFeedLength int
	webFeedLength int
	// conf
	c *conf.Config
	// bangumi http client
	httpClient *bm.Client
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// redis
		redis:                  redis.NewPool(c.MultiRedis.Cache),
		redisTTLUpper:          int32(time.Duration(c.MultiRedis.TTLUpper) / time.Second),
		redisExpireUpper:       int32(time.Duration(c.MultiRedis.ExpireUpper) / time.Second),
		redisExpireFeed:        int32(time.Duration(c.MultiRedis.ExpireFeed) / time.Second),
		redisExpireArchiveFeed: int32(time.Duration(c.Feed.ArchiveFeedExpire) / time.Second),
		redisExpireBangumiFeed: int32(time.Duration(c.Feed.BangumiFeedExpire) / time.Second),
		// mc
		mc:            memcache.NewPool(c.Memcache.Config),
		mcExpire:      int32(time.Duration(c.Memcache.Expire) / time.Second),
		bangumiExpire: int32(time.Duration(c.Memcache.Expire) / time.Second),
		// feed Config
		appFeedLength: c.Feed.AppLength,
		webFeedLength: c.Feed.WebLength,
		httpClient:    bm.NewClient(c.HTTPClient),
	}
	if d.appFeedLength == 0 {
		d.appFeedLength = 200
	}
	if d.webFeedLength == 0 {
		d.webFeedLength = 400
	}
	return
}

// Ping ping health of redis and mc.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingRedis(c); err != nil {
		return
	}
	return d.pingMC(c)
}

// Close close connections of redis and mc.
func (d *Dao) Close() {
	if d.redis != nil {
		d.redis.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
}
