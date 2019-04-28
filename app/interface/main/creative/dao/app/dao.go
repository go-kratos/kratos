package app

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/library/database/sql"
)

// Dao  define
type Dao struct {
	db *sql.DB
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.DB.Creative),
	}
	return
}

// Ping db
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close db
func (d *Dao) Close() (err error) {
	return d.db.Close()
}
