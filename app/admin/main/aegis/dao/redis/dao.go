package redis

import (
	"context"

	"go-common/app/admin/main/aegis/conf"
	"go-common/library/cache/redis"
)

// Dao dao
type Dao struct {
	c         *conf.Config
	cluster   *redis.Pool
	netExpire int64
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:         c,
		cluster:   redis.NewPool(c.Redis.Cluster),
		netExpire: int64(c.Redis.NetExpire),
	}

	if dao.netExpire == 0 {
		dao.netExpire = int64(1800) //30m
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.cluster.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}
