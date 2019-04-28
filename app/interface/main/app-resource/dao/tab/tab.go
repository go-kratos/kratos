package tab

import (
	"context"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model/tab"
	"go-common/library/database/sql"
)

const (
	_getAllMenuSQL = "SELECT id,plat,name,ctype,cvalue,plat_ver,status,color,badge FROM app_menus WHERE stime<? AND etime>? AND status=1 ORDER BY `order` ASC"
)

type Dao struct {
	db *sql.DB
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	return
}

// Menus menus tab
func (d *Dao) Menus(c context.Context, now time.Time) (menus []*tab.Menu, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getAllMenuSQL, now, now); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &tab.Menu{}
		if err = rows.Scan(&m.TabID, &m.Plat, &m.Name, &m.CType, &m.CValue, &m.PlatVersion, &m.Status, &m.Color, &m.Badge); err != nil {
			return
		}
		if m.CValue != "" {
			m.Change()
			menus = append(menus, m)
		}
	}
	err = rows.Err()
	return
}
