package dao

import (
	"context"

	"go-common/app/admin/main/point/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/elastic"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

const _searchBussinss = "vip_point_change_history"

// Dao dao
type Dao struct {
	c      *conf.Config
	mc     *memcache.Pool
	db     *xsql.DB
	client *bm.Client
	es     *elastic.Elastic
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:      c,
		mc:     memcache.NewPool(c.Memcache),
		db:     xsql.NewMySQL(c.MySQL),
		client: bm.NewClient(c.HTTPClient),
		// es
		es: elastic.NewElastic(nil),
	}
	return
}

// Close close the resource.
func (dao *Dao) Close() {
	dao.mc.Close()
	dao.db.Close()
}

// Ping dao ping
func (dao *Dao) Ping(c context.Context) error {
	return dao.pingMC(c)
}

// pingMc ping
func (dao *Dao) pingMC(c context.Context) (err error) {
	conn := dao.mc.Get(c)
	defer conn.Close()
	return
}
