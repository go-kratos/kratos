package task

import (
	"go-common/app/admin/main/videoup/conf"
	"go-common/library/cache/redis"
)

// Dao is track dao.
type Dao struct {
	c *conf.Config
	// redis
	redis *redis.Pool
}

var (
	d *Dao
)

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:     c,
		redis: redis.NewPool(c.Redis.Secondary.Config),
	}
	return d
}
