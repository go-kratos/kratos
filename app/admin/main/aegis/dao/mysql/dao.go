package mysql

import (
	"context"

	"go-common/app/admin/main/aegis/conf"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c  *conf.Config
	db *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c: c,

		db: xsql.NewMySQL(c.MySQL),
	}

	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return d.db.Ping(c)
}
