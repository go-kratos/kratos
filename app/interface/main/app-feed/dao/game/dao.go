package game

import (
	"context"

	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-feed/conf"
	"go-common/library/database/sql"
)

const (
	_getSQL = "SELECT `id`,`title`,`desc`,`icon`,`cover`,`url_type`,`url_value`,`btn_txt`,`re_type`,`re_value`,`number`,`double_cover` FROM download_card"
)

type Dao struct {
	db          *sql.DB
	downloadGet *sql.Stmt
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Manager),
	}
	// prepare
	d.downloadGet = d.db.Prepared(_getSQL)
	return
}

func (d *Dao) DownLoad(c context.Context) (dm map[int64]*operate.Download, err error) {
	rows, err := d.downloadGet.Query(c)
	if err != nil {
		return
	}
	defer rows.Close()
	dm = map[int64]*operate.Download{}
	for rows.Next() {
		d := &operate.Download{}
		if err = rows.Scan(&d.ID, &d.Title, &d.Desc, &d.Icon, &d.Cover, &d.URLType, &d.URLValue, &d.BtnTxt, &d.ReType, &d.ReValue, &d.Number, &d.DoubleCover); err != nil {
			return
		}
		d.Change()
		dm[d.ID] = d
	}
	return
}

// Close close memcache resource.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
