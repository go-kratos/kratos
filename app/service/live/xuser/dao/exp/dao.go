package exp

import (
	"context"
	"go-common/app/service/live/xuser/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
)

// Dao exp dao
type Dao struct {
	c        *conf.Config
	db       *xsql.DB
	memcache *memcache.Pool
}

// NewExpDao init mysql db
func NewExpDao(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:        c,
		db:       xsql.NewMySQL(c.UserExpMySQL),
		memcache: memcache.NewPool(c.ExpMemcache),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
	d.memcache.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return nil
}

func (d *Dao) getExpire() (respExpire int32) {
	if t := conf.Conf.UserExpExpire; t != nil {
		respExpire = t.ExpireTime
	} else {
		respExpire = _emptyExpire
	}
	return
}
