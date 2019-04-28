package archive

import (
	"time"

	"go-common/app/job/main/tv/conf"
	"go-common/library/cache/memcache"
)

// Dao is archive dao.
type Dao struct {
	// memcache
	mc         *memcache.Pool
	expireView int32
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// memcache
		mc:         memcache.NewPool(c.Memcache.Config),
		expireView: int32(time.Duration(c.Memcache.ExpireMedia) / time.Second),
	}
	return
}
