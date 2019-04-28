package dao

import (
	"go-common/app/job/main/identify/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c      *conf.Config
	authDB *xsql.DB
	authMC *memcache.Pool
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:      c,
		authDB: xsql.NewMySQL(c.AuthDB),
		authMC: memcache.NewPool(c.AuthMC),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.authDB.Close()
}
