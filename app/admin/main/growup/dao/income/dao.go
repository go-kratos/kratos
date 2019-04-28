package income

import (
	"context"

	"go-common/app/admin/main/growup/conf"
	"go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c  *conf.Config
	db *sql.DB
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.ORM.Allowance),
	}
	return d
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close close connections of mc, redis, db.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// BeginTran begin transcation
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}
