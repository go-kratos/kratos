package dao

import (
	"go-common/app/job/main/click/conf"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao is
type Dao struct {
	c      *conf.Config
	db     *sql.DB
	client *bm.Client
}

// New is
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		db:     sql.NewMySQL(c.DB),
		client: bm.NewClient(c.HTTPClient),
	}
	return d
}
