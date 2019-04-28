package pgc

import (
	"go-common/app/interface/main/tv/conf"
	"go-common/library/cache/memcache"
	bm "go-common/library/net/http/blademaster"
)

// Dao is account dao.
type Dao struct {
	conf   *conf.Config
	client *bm.Client
	mc     *memcache.Pool
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf:   c,
		client: bm.NewClient(c.HTTPClient),
		mc:     memcache.NewPool(c.Memcache.Config),
	}
	return
}
