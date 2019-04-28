package dao

import (
	"context"

	"go-common/app/service/main/sms/conf"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

// Dao struct info of Dao.
type Dao struct {
	c       *conf.Config
	db      *xsql.DB
	client  *bm.Client
	databus *databus.Databus
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:       c,
		db:      xsql.NewMySQL(c.MySQL),
		client:  bm.NewClient(c.HTTPClient),
		databus: databus.New(c.Databus),
	}
	return
}

// Ping ping health of db.
func (d *Dao) Ping(ctx context.Context) (err error) {
	return d.db.Ping(ctx)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	d.db.Close()
	d.databus.Close()
}
