package dao

import (
	"context"

	"go-common/app/job/main/ugcpay/conf"
	ugcpay_rank "go-common/app/service/main/ugcpay-rank/api/v1"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c             *conf.Config
	mc            *memcache.Pool
	mcRank        *memcache.Pool
	redis         *redis.Pool
	db            *xsql.DB
	dbrank        *xsql.DB
	dbrankold     *xsql.DB
	ugcPayRankAPI ugcpay_rank.UGCPayRankClient
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:         c,
		mc:        memcache.NewPool(c.Memcache),
		mcRank:    memcache.NewPool(c.MemcacheRank),
		redis:     redis.NewPool(c.Redis),
		db:        xsql.NewMySQL(c.MySQL),
		dbrank:    xsql.NewMySQL(c.MySQLRank),
		dbrankold: xsql.NewMySQL(c.MySQLRankOld),
	}
	var err error
	if dao.ugcPayRankAPI, err = ugcpay_rank.NewClient(c.GRPCUGCPayRank); err != nil {
		panic(err)
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*xsql.Tx, error) {
	return d.db.Begin(c)
}

// BeginTranRank .
func (d *Dao) BeginTranRank(c context.Context) (*xsql.Tx, error) {
	return d.dbrank.Begin(c)
}
