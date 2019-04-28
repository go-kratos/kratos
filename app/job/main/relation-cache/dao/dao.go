package dao

import (
	"context"
	"time"

	"go-common/app/job/main/relation-cache/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

// Dao dao
type Dao struct {
	*cacheTTL
	c      *conf.Config
	mc     *memcache.Pool
	redis  *redis.Pool
	db     *xsql.DB
	client *bm.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		cacheTTL: &cacheTTL{
			RelationTTL: asSecond(c.CacheTTL.RelationTTL),
		},
		c:      c,
		mc:     memcache.NewPool(c.Memcache),
		redis:  redis.NewPool(c.Redis),
		db:     xsql.NewMySQL(c.MySQL),
		client: bm.NewClient(c.HTTPClient),
	}
	return
}

type cacheTTL struct {
	RelationTTL int64
}

func asSecond(d xtime.Duration) int64 {
	return int64(time.Duration(d) / time.Second)
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}
