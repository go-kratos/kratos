package charge

import (
	"context"

	"go-common/app/job/main/growup/conf"
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
		db: sql.NewMySQL(c.Mysql.Growup),
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

// Exec exec
func (d *Dao) Exec(c context.Context, sql string) (err error) {
	_, err = d.db.Exec(c, sql)
	return
}

// QueryRow QueryRow
func (d *Dao) QueryRow(c context.Context, sql string) (rows *sql.Row) {
	return d.db.QueryRow(c, sql)
}

// Query query
func (d *Dao) Query(c context.Context, sql string) (rows *sql.Rows, err error) {
	return d.db.Query(c, sql)
}
