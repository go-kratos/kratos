package dao

import (
	"context"

	"go-common/app/infra/notify/conf"
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
		c:  c,
		db: xsql.NewMySQL(c.MySQL),
	}
	return
}

// Close close the resource.
func (dao *Dao) Close() {
	dao.db.Close()
}

// Ping dao ping
func (dao *Dao) Ping(c context.Context) error {
	return dao.db.Ping(c)
}
