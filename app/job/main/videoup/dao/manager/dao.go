package manager

import (
	"go-common/app/job/main/videoup/conf"
	xsql "go-common/library/database/sql"
)

// Dao is manager dao.
type Dao struct {
	c  *conf.Config
	db *xsql.DB
}

// New new a manager dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: xsql.NewMySQL(c.DB.Manager),
	}
	return d
}
