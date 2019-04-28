package dao

import (
	"context"

	"go-common/app/job/main/stat/conf"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao is stat job dao.
type Dao struct {
	c         *conf.Config
	smsClient *bm.Client
	db        *xsql.DB
	clickDB   *xsql.DB
}

// New add a feed job dao.
func New(c *conf.Config) *Dao {
	return &Dao{
		c:         c,
		smsClient: bm.NewClient(c.HTTPClient),
		db:        xsql.NewMySQL(c.DB),
		clickDB:   xsql.NewMySQL(c.ClickDB),
	}
}

// Ping ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}
