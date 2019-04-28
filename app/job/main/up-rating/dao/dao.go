package dao

import (
	"context"

	"go-common/app/job/main/up-rating/conf"

	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao dao
type Dao struct {
	c  *conf.Config
	db *sql.DB
}

// New fn
func New(c *conf.Config) (d *Dao) {
	log.Info("dao start")
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.MySQL.Rating),
	}
	//d.db.State = prom.LibClient
	return
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// BeginTran begin transcation
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}
