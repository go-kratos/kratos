package special

import (
	"context"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/special"
	"go-common/library/database/sql"
)

const (
	_getSQL = "SELECT `id`,`title`,`desc`,`cover`,`re_type`,`re_value`,`corner` FROM special_card"
)

type Dao struct {
	db         *sql.DB
	specialGet *sql.Stmt
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Manager),
	}
	// prepare
	d.specialGet = d.db.Prepared(_getSQL)
	return
}

func (d *Dao) Card(c context.Context) (scm map[int64]*special.Card, err error) {
	rows, err := d.specialGet.Query(c)
	if err != nil {
		return
	}
	defer rows.Close()
	scm = map[int64]*special.Card{}
	for rows.Next() {
		sc := &special.Card{}
		if err = rows.Scan(&sc.ID, &sc.Title, &sc.Desc, &sc.Cover, &sc.ReType, &sc.ReValue, &sc.Badge); err != nil {
			return
		}
		scm[sc.ID] = sc
	}
	return scm, err
}

// Close close memcache resource.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
