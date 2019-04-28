package dao

import (
	"time"

	"go-common/app/interface/main/tv/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/stat/prom"
)

// Dao .
type Dao struct {
	db         *sql.DB
	conf       *conf.Config
	client     *bm.Client
	redis      *redis.Pool
	mc         *memcache.Pool
	dbeiExpire int64
}

// New .
func New(c *conf.Config) *Dao {
	return &Dao{
		db:         sql.NewMySQL(c.Mysql),
		conf:       c,
		client:     bm.NewClient(c.HTTPClient),
		redis:      redis.NewPool(c.Redis.Config),
		dbeiExpire: int64(time.Duration(c.Cfg.Dangbei.Expire) / time.Second),
		mc:         memcache.NewPool(c.Memcache.Config),
	}
}

// Prom
var (
	errorsCount = prom.BusinessErrCount
	infosCount  = prom.BusinessInfoCount
)

// PromError prom error
func PromError(name string) {
	errorsCount.Incr(name)
}

// PromInfo add prom info
func PromInfo(name string) {
	infosCount.Incr(name)
}
