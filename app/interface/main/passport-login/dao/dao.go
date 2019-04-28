package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/passport-login/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c            *conf.Config
	authMC       *memcache.Pool
	authMCExpire int32
	originDB     *xsql.DB
	userDB       *xsql.DB
	secretDB     *xsql.DB
	authDB       *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:            c,
		authMC:       memcache.NewPool(c.Memcache.Auth.Config),
		authMCExpire: int32(time.Duration(c.Memcache.Auth.Expire) / time.Second),
		originDB:     xsql.NewMySQL(c.DB.Origin),
		userDB:       xsql.NewMySQL(c.DB.User),
		secretDB:     xsql.NewMySQL(c.DB.Secret),
		authDB:       xsql.NewMySQL(c.DB.Auth),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.authMC.Close()
	d.originDB.Close()
	d.userDB.Close()
	d.secretDB.Close()
	d.authDB.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
