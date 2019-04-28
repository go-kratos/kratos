package manager

import (
	"context"

	"go-common/app/service/main/videoup/conf"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

const (
	_searchUpdateURL = "/x/admin/search/upsert"
)

// Dao is redis dao.
type Dao struct {
	c         *conf.Config
	managerDB *sql.DB
	// select stmt
	upsStmt         *sql.Stmt
	httpW           *bm.Client
	searchUpdateURL string
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:               c,
		managerDB:       sql.NewMySQL(c.DB.Manager),
		httpW:           bm.NewClient(c.HTTPClient.Write),
		searchUpdateURL: c.Host.Manager + _searchUpdateURL,
	}
	// select stmt
	d.upsStmt = d.managerDB.Prepared(_upsSQL)
	return d
}

// Close fn
func (d *Dao) Close() {
	if d.managerDB != nil {
		d.managerDB.Close()
	}
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return d.managerDB.Ping(c)
}
