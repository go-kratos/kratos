package audit

import (
	"context"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_auditSQL = "SELECT mobi_app,build FROM audit"
)

// Dao is notice dao.
type Dao struct {
	db    *sql.DB
	audit *sql.Stmt
}

// New new a notice dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	d.audit = d.db.Prepared(_auditSQL)
	return
}

// Audits get all audit build.
func (d *Dao) Audits(ctx context.Context) (res map[string]map[int]struct{}, err error) {
	rows, err := d.audit.Query(ctx)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	defer rows.Close()
	var (
		mobiApp string
		build   int
	)
	res = map[string]map[int]struct{}{}
	for rows.Next() {
		if err = rows.Scan(&mobiApp, &build); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		if plat, ok := res[mobiApp]; ok {
			plat[build] = struct{}{}
		} else {
			res[mobiApp] = map[int]struct{}{
				build: struct{}{},
			}
		}
	}
	return
}

// Close close memcache resource.
func (dao *Dao) Close() {
	if dao.db != nil {
		dao.db.Close()
	}
}
