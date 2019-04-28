package redis

import (
	"go-common/app/job/main/videoup-report/conf"
	"go-common/library/cache/redis"
)

// Dao is redis dao.
type Dao struct {
	c         *conf.Config
	redis     *redis.Pool
	secondary *redis.Pool
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:         c,
		redis:     redis.NewPool(c.Redis.Track.Config),
		secondary: redis.NewPool(c.Redis.Secondary.Config),
	}
	return d
}

// Close close the redis connection
func (d *Dao) Close() (err error) {
	if err = d.secondary.Close(); err != nil {
		return
	}
	return d.redis.Close()
}
