package toview

import (
	"context"
	"time"

	"go-common/app/interface/main/history/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/hbase.v2"
)

// Dao dao.
type Dao struct {
	conf   *conf.Config
	info   *hbase.Client
	redis  *redis.Pool
	expire int
}

// New new history dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:   c,
		info:   hbase.NewClient(c.Info.Config),
		redis:  redis.NewPool(c.Toview.Config),
		expire: int(time.Duration(c.Toview.Expire) / time.Second),
	}
	return
}

// Ping check connection success.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.PingRedis(c)
}

// Close close the redis and kafka resource.
func (d *Dao) Close() {
	if d.redis != nil {
		d.redis.Close()
	}
}
