package watermark

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/library/database/sql"
	httpx "go-common/library/net/http/blademaster"
)

// Dao  define
type Dao struct {
	c  *conf.Config
	db *sql.DB
	// http client
	client *httpx.Client
	genWm  string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.DB.Creative),
		// http client
		client: httpx.NewClient(c.HTTPClient.Normal),
		genWm:  c.Host.API + _genWm,
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
