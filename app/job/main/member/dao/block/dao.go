package block

import (
	"context"

	"go-common/app/job/main/member/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

type notifyFunc func(context.Context, int64, string) error

// Dao dao
type Dao struct {
	conf       *conf.Config
	mc         *memcache.Pool
	db         *xsql.DB
	httpClient *bm.Client
	notifyFunc notifyFunc
}

// New init mysql db
func New(conf *conf.Config, mc *memcache.Pool, db *xsql.DB, client *bm.Client, notifyFunc notifyFunc) (dao *Dao) {
	dao = &Dao{
		conf:       conf,
		mc:         mc,
		db:         db,
		httpClient: client,
		notifyFunc: notifyFunc,
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

// BeginTX is.
func (d *Dao) BeginTX(c context.Context) (tx *xsql.Tx, err error) {
	if tx, err = d.db.Begin(c); err != nil {
		err = errors.WithStack(err)
	}
	return
}
