package dao

import (
	"context"

	"go-common/app/job/main/credit-timer/conf"
	"go-common/library/database/sql"
)

// Dao struct info of Dao.
type Dao struct {
	db *sql.DB
	c  *conf.Config
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.Mysql),
	}
	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// Ping ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}
