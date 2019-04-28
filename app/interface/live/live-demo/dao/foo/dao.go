package foo

import (
	"context"

	"go-common/app/interface/live/live-demo/conf"
	"go-common/library/cache/redis"
)

// Dao dao
type Dao struct {
	c     *conf.Config
	redis *redis.Pool
}

// New init
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		redis: redis.NewPool(c.Redis),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	_, err := d.redis.Get(c).Do("ping")
	return err
}
