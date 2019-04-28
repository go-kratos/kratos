package show

import (
	"context"

	"go-common/app/service/main/resource/conf"
	xsql "go-common/library/database/sql"
)

// Dao is resource dao.
type Dao struct {
	db *xsql.DB
	c  *conf.Config
}

// New init mysql db
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: xsql.NewMySQL(c.DB.Show),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
}

// Ping check dao health.
func (d *Dao) Ping(c context.Context) error {
	return d.db.Ping(c)
}
