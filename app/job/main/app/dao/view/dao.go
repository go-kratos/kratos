package view

import (
	"context"
	"time"

	"go-common/app/job/main/app/conf"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"

	httpx "go-common/library/net/http/blademaster"
)

const (
	_configHost = "/v1/config/host/infos"
)

// Dao is dao.
type Dao struct {
	client *httpx.Client
	config string
	// mc
	mc       *memcache.Pool
	expireMc int32
	// redis
	redis *redis.Pool
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPClient),
		config: c.Host.Config + _configHost,
		// mc
		mc:       memcache.NewPool(c.Memcache.Feed.Config),
		expireMc: int32(time.Duration(c.Memcache.Feed.ExpireMaxAid) / time.Second),
		// reids
		redis: redis.NewPool(c.Redis.Feed.Config),
	}
	return
}

func (d *Dao) PingMc(c context.Context) (err error) {
	conn := d.mc.Get(c)
	item := &memcache.Item{Key: "ping", Value: []byte{1}, Flags: memcache.FlagRAW, Expiration: 0}
	err = conn.Set(item)
	conn.Close()
	return
}
