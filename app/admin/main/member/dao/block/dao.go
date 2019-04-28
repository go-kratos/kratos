package block

import (
	"context"

	"go-common/app/admin/main/member/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// Dao .
type Dao struct {
	conf       *conf.Config
	mc         *memcache.Pool
	db         *xsql.DB
	httpClient *bm.Client
}

// New init mysql db
func New(conf *conf.Config, client *bm.Client, mc *memcache.Pool, db *xsql.DB) (dao *Dao) {
	dao = &Dao{
		conf:       conf,
		mc:         mc,
		db:         db,
		httpClient: client,
	}
	return
}

// BeginTX .
func (d *Dao) BeginTX(c context.Context) (tx *xsql.Tx, err error) {
	if tx, err = d.db.Begin(c); err != nil {
		err = errors.WithStack(err)
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
	return
}
