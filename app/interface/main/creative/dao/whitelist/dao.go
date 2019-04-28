package whitelist

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

var (
	_getAllSQL = "select mid, type from whitelist where state=1 order by type asc, ctime desc"
)

// Dao  define
type Dao struct {
	db         *sql.DB
	getAllStmt *sql.Stmt
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.DB.Creative),
	}
	d.getAllStmt = d.db.Prepared(_getAllSQL)
	return
}

// Ping db
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}

// Close db
func (d *Dao) Close() (err error) {
	return d.db.Close()
}

// List fn
func (d *Dao) List(c context.Context) (wls []*archive.WhiteList, err error) {
	rows, err := d.getAllStmt.Query(c)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		wl := &archive.WhiteList{}
		if err = rows.Scan(&wl.Mid, &wl.Tp); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		wls = append(wls, wl)
	}
	return
}
