package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/spy/conf"
	account "go-common/app/service/main/account/rpc/client"
	spy "go-common/app/service/main/spy/rpc/client"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
)

// Dao struct user of Dao.
type Dao struct {
	c *conf.Config
	// db
	db              *sql.DB
	getUserInfoStmt []*sql.Stmt
	eventStmt       *sql.Stmt
	factorAllStmt   *sql.Stmt
	allGroupStmt    *sql.Stmt
	// cache
	mcUser       *memcache.Pool
	mcUserExpire int32
	// rpc
	spyRPC *spy.Service
	accRPC *account.Service3
}

// New create a instance of Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// db
		db: sql.NewMySQL(c.DB.Spy),
		// mc
		mcUser:       memcache.NewPool(c.Memcache.User),
		mcUserExpire: int32(time.Duration(c.Memcache.UserExpire) / time.Second),
		// rpc
		spyRPC: spy.New(c.SpyRPC),
		accRPC: account.New3(c.AccountRPC),
	}
	if conf.Conf.Property.UserInfoShard <= 0 {
		panic("conf.Conf.Property.UserInfoShard <= 0")
	}
	if conf.Conf.Property.HistoryShard <= 0 {
		panic("conf.Conf.Property.HistoryShard <= 0")
	}
	d.getUserInfoStmt = make([]*sql.Stmt, conf.Conf.Property.UserInfoShard)
	for i := int64(0); i < conf.Conf.Property.UserInfoShard; i++ {
		d.getUserInfoStmt[i] = d.db.Prepared(fmt.Sprintf(_getUserInfoSQL, i))
	}
	d.eventStmt = d.db.Prepared(_eventSQL)
	d.factorAllStmt = d.db.Prepared(_factorAllSQL)
	d.allGroupStmt = d.db.Prepared(_allGroupSQL)
	return
}

// Ping check connection of db , mc.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close close connection of db , mc.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
