package caldiff

import (
	"go-common/app/job/main/appstatic/conf"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

// Dao .
type Dao struct {
	c      *conf.Config
	db     *xsql.DB
	client *bm.Client
}

// New creates a dao instance.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		db:     xsql.NewMySQL(c.MySQL),
		client: bm.NewClient(c.HTTPClient),
	}
	return
}
