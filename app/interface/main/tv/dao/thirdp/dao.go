package thirdp

import (
	"runtime"
	"time"

	"go-common/app/interface/main/tv/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

// Dao .
type Dao struct {
	db         *sql.DB
	conf       *conf.Config
	redis      *redis.Pool
	mc         *memcache.Pool
	dbeiExpire int64
	cntExpire  int32
	mCh        chan func()
}

// New .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:         sql.NewMySQL(c.Mysql),
		conf:       c,
		redis:      redis.NewPool(c.Redis.Config),
		dbeiExpire: int64(time.Duration(c.Cfg.Dangbei.Expire) / time.Second),
		mc:         memcache.NewPool(c.Memcache.Config),
		cntExpire:  int32(time.Duration(c.Memcache.MangoExpire) / time.Second),
		mCh:        make(chan func(), 10240),
	}
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go d.cacheproc()
	}
	return
}

var (
	cachedCount = prom.CacheHit
	missedCount = prom.CacheMiss
)

// addCache add archive to mc or redis
func (d *Dao) addCache(f func()) {
	select {
	case d.mCh <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc write memcache and stat redis use goroutine
func (d *Dao) cacheproc() {
	for {
		f := <-d.mCh
		f()
	}
}
