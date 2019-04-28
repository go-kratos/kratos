package dao

import (
	"context"
	"net/http"
	"time"

	"go-common/app/job/main/answer/conf"
	"go-common/library/database/sql"
	xhttp "go-common/library/net/http/blademaster"
)

const (
	_bfsTimeout = 5 * time.Second
	_beFormal   = "/api/internal/member/beFormal"
)

// Dao event dao def.
type Dao struct {
	c        *conf.Config
	db       *sql.DB
	client   *http.Client
	xclient  *xhttp.Client
	beFormal string
}

// New create instance of dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.Mysql),
		client: &http.Client{
			Timeout: _bfsTimeout,
		},
		xclient:  xhttp.NewClient(c.HTTPClient),
		beFormal: c.Properties.AccountIntranetURI + _beFormal,
	}
	return
}

// Ping check db health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close close all db connections.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
