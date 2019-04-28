package special

import (
	"context"
	"time"

	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-feed/conf"
	"go-common/library/database/sql"
)

const (
	_getSQL = "SELECT `id`,`title`,`desc`,`cover`,`scover`,`re_type`,`re_value`,`corner`,`size` FROM special_card"
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

func (d *Dao) Card(c context.Context, now time.Time) (scm map[int64]*operate.Special, err error) {
	rows, err := d.specialGet.Query(c)
	if err != nil {
		return
	}
	defer rows.Close()
	scm = map[int64]*operate.Special{}
	for rows.Next() {
		sc := &operate.Special{}
		if err = rows.Scan(&sc.ID, &sc.Title, &sc.Desc, &sc.Cover, &sc.SingleCover, &sc.ReType, &sc.ReValue, &sc.Badge, &sc.Size); err != nil {
			return
		}
		sc.Change()
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
