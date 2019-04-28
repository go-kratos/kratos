package dao

import (
	"context"

	"go-common/app/job/main/block/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// Dao dao
type Dao struct {
	mc         *memcache.Pool
	db         *xsql.DB
	httpClient *bm.Client
}

// New init mysql db
func New() (dao *Dao) {
	dao = &Dao{
		mc:         memcache.NewPool(conf.Conf.Memcache),
		db:         xsql.NewMySQL(conf.Conf.DB),
		httpClient: bm.NewClient(conf.Conf.HTTPClient),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		return
	}
	if err = d.pingMC(c); err != nil {
		return
	}
	return
}

// pingMc ping
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	return
}

func (d *Dao) BeginTX(c context.Context) (tx *xsql.Tx, err error) {
	if tx, err = d.db.Begin(c); err != nil {
		err = errors.WithStack(err)
	}
	return
}
