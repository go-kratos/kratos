package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/tag/conf"
	"go-common/app/admin/main/tag/model"
	accwarden "go-common/app/service/main/account/api"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

const (
	_shard = 200
)

// Dao dao layer.
type Dao struct {
	es            *elastic.Elastic
	db            *sql.DB
	redis         *redis.Pool
	redisRank     *redis.Pool
	mc            *memcache.Pool
	client        *bm.Client
	hosts         *model.DependServiceHost
	accClient     accwarden.AccountClient
	channelExpire int32
}

// New new a dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		es:            elastic.NewElastic(c.ES),
		db:            sql.NewMySQL(c.Mysql),
		redis:         redis.NewPool(c.Redis.Tag),
		redisRank:     redis.NewPool(c.Redis.TagRank),
		mc:            memcache.NewPool(c.Memcache.Tag),
		client:        bm.NewClient(c.HTTPClient),
		channelExpire: int32(time.Duration(c.Memcache.ChannelExpire) / time.Second),
		hosts:         c.Hosts,
	}
	var err error
	if d.accClient, err = accwarden.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	return
}

func (d *Dao) hit(id int64) string {
	return fmt.Sprintf("%03d", id%int64(_shard))
}

// Ping check health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close when stoped, relese all resource.
func (d *Dao) Close(c context.Context) (err error) {
	return d.db.Close()
}

// BeginTran begin mysql transaction
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}
