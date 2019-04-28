package income

import (
	"context"

	"go-common/app/job/main/growup/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao dao
type Dao struct {
	c    *conf.Config
	db   *sql.DB
	rddb *sql.DB
	Tx   *sql.Tx
}

// New fn
func New(c *conf.Config) (d *Dao) {
	log.Info("dao start")
	d = &Dao{
		c:    c,
		db:   sql.NewMySQL(c.Mysql.Growup),
		rddb: sql.NewMySQL(c.Mysql.Allowance),
	}
	d.Tx, _ = d.BeginTran(context.TODO())
	//d.db.State = prom.LibClient
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

// Exec exec sql
func (d *Dao) Exec(c context.Context, sql string) (err error) {
	_, err = d.db.Exec(c, sql)
	return
}

// QueryRow query row
func (d *Dao) QueryRow(c context.Context, sql string) (rows *sql.Row) {
	return d.db.QueryRow(c, sql)
}

// Query query
func (d *Dao) Query(c context.Context, sql string) (rows *sql.Rows, err error) {
	return d.db.Query(c, sql)
}

// test
// func (d *Dao) Truncate(c context.Context, table string) {
// 	d.db.Exec(c, fmt.Sprintf("truncate %s", table))
// }
