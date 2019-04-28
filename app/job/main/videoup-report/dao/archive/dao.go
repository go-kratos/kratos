package archive

import (
	"context"

	"go-common/app/job/main/videoup-report/conf"
	"go-common/library/database/sql"
)

// Dao is redis dao.
type Dao struct {
	c *conf.Config
	// db
	db *sql.DB
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.DB.Archive),
	}
	// select
	return d
}

// BeginTran begin transcation.
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.db.Begin(c)
}

// Close close dao.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}
