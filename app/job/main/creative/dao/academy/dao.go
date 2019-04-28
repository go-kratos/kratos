package academy

import (
	"context"
	"go-common/app/job/main/creative/conf"
	"go-common/library/database/sql"
)

// Dao is creative dao.
type Dao struct {
	// config
	c *conf.Config
	// db
	db *sql.DB
}

// New init api url
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.DB.Creative),
	}
	return
}

// Ping creativeDb
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close creativeDb
func (d *Dao) Close() (err error) {
	return d.db.Close()
}
