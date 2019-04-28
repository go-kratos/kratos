package up

import (
	"context"

	"go-common/app/admin/main/mcn/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao
type Dao struct {
	c                   *conf.Config
	mc                  *memcache.Pool
	db                  *xsql.DB
	client              *bm.Client
	arcTopURL           string
	dataFansURL         string
	dataFansBaseAttrURL string
	dataFansAreaURL     string
	dataFansTypeURL     string
	dataFansTagURL      string
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:  c,
		mc: memcache.NewPool(c.Memcache),
		db: xsql.NewMySQL(c.MySQL),
		// http client
		client: bm.NewClient(c.HTTPClient),
		// url
		arcTopURL:           c.Host.API + _arcTopURL,
		dataFansURL:         c.Host.API + _dataFansURL,
		dataFansBaseAttrURL: c.Host.API + _dataFansBaseAttrURL,
		dataFansAreaURL:     c.Host.API + _dataFansAreaURL,
		dataFansTypeURL:     c.Host.API + _dataFansTypeURL,
		dataFansTagURL:      c.Host.API + _dataFansTagURL,
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.db.Ping(c)
}

// BeginTran start tx .
func (d *Dao) BeginTran(c context.Context) (tx *xsql.Tx, err error) {
	return d.db.Begin(c)
}
