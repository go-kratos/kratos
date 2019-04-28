package redis

import (
	"go-common/app/job/main/videoup/conf"
	xredis "go-common/library/cache/redis"
)

// Dao is redis dao.
type Dao struct {
	c     *conf.Config
	redis *xredis.Pool
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:     c,
		redis: xredis.NewPool(c.Redis),
	}
	return d
}

// Close close the redis connection
func (d *Dao) Close() (err error) {
	return d.redis.Close()
}
