package academy

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
)

// Dao  define
type Dao struct {
	c  *conf.Config
	db *sql.DB
	es *elastic.Elastic
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.DB.Creative),
		es: elastic.NewElastic(nil),
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
