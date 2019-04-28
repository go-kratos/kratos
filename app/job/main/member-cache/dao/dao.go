package dao

import (
	"context"

	"go-common/app/job/main/member-cache/conf"
	"go-common/library/cache/memcache"
)

// Dao dao
type Dao struct {
	c              *conf.Config
	memberMemcache *memcache.Pool
	blockMemcache  *memcache.Pool
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:              c,
		memberMemcache: memcache.NewPool(c.MemberMemcache),
		blockMemcache:  memcache.NewPool(c.BlockMemcache),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.memberMemcache.Close()
	d.blockMemcache.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}
