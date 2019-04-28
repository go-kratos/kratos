package monitor

import (
	"go-common/app/admin/main/videoup/conf"
	"go-common/library/cache/redis"
)

// Dao is redis dao.
type Dao struct {
	c     *conf.Config
	redis *redis.Pool
}

var (
	d *Dao
)

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:     c,
		redis: redis.NewPool(c.Redis.Secondary.Config),
	}
	return d
}

// Close close dao.
func (d *Dao) Close() {
	d.redis.Close()
}
