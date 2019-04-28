package share

import (
	"context"

	"go-common/library/cache/redis"
	"go-common/library/queue/databus"

	"go-common/app/service/main/archive/conf"
)

// Dao is share dao.
type Dao struct {
	c *conf.Config
	// redis
	rds    *redis.Pool
	expire int32
	// databus
	statDbus  *databus.Databus
	shareDbus *databus.Databus
}

// New new a share dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:         c,
		rds:       redis.NewPool(c.Redis.Archive.Config),
		expire:    168 * 60 * 60,
		statDbus:  databus.New(c.StatDatabus),
		shareDbus: databus.New(c.ShareDatabus),
	}
	return d
}

// Close close resource.
func (d *Dao) Close() {
	if d.rds != nil {
		d.rds.Close()
	}
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	conn := d.rds.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}
