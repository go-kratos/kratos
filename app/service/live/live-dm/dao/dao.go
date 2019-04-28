package dao

import (
	"context"

	"go-common/app/service/live/live-dm/conf"
	"go-common/library/cache/redis"
	"go-common/library/sync/pipeline/fanout"
)

// Dao dao
type Dao struct {
	c              *conf.Config
	redis          *redis.Pool
	whitelistredis *redis.Pool
	Databus        *fanout.Fanout
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:              c,
		redis:          redis.NewPool(c.Redis),
		whitelistredis: redis.NewPool(c.WhiteListRedis),
		Databus:        fanout.New("dmDatabus", fanout.Worker(1), fanout.Buffer(c.CacheDatabus.Size)),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.Databus.Close()
	d.whitelistredis.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: redis
	return nil
}
