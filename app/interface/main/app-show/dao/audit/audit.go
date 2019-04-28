package audit

import (
	"context"

	"go-common/app/interface/main/app-show/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getSQL = "SELECT mobi_app,build FROM audit"
)

// Dao is audit dao.
type Dao struct {
	db     *sql.DB
	audGet *sql.Stmt
}

// New new a audit dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	d.audGet = d.db.Prepared(_getSQL)
	return
}

// Audits get all audit build.
func (d *Dao) Audits(ctx context.Context) (res map[string]map[int]struct{}, err error) {
	rows, err := d.audGet.Query(ctx)
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
