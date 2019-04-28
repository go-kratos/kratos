package material

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/library/database/sql"
)

// Dao is archive dao.
type Dao struct {
	// config
	c  *conf.Config
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

// Ping fn
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close fn
func (d *Dao) Close() (err error) {
	return d.db.Close()
}
