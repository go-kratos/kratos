package tab

import (
	"context"

	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-feed/conf"

	"go-common/library/database/sql"
)

const (
	_getAllMenuSQL   = "SELECT id,plat,name,ctype,cvalue,plat_ver,stime,etime,status,color,badge FROM app_menus ORDER BY `order`"
	_getAllActiveSQL = "SELECT id,parent_id,name,background,type,content FROM app_active"
)

type Dao struct {
	db        *sql.DB
	menuGet   *sql.Stmt
	activeGet *sql.Stmt
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	// prepare
	d.menuGet = d.db.Prepared(_getAllMenuSQL)
	d.activeGet = d.db.Prepared(_getAllActiveSQL)
	return
}

func (d *Dao) Menus(c context.Context) (menus []*operate.Menu, err error) {
	var rows *sql.Rows
	if rows, err = d.menuGet.Query(c); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &operate.Menu{}
		if err = rows.Scan(&m.TabID, &m.Plat, &m.Name, &m.CType, &m.CValue, &m.PlatVersion, &m.STime, &m.ETime, &m.Status, &m.Color, &m.Badge); err != nil {
			return
		}
		if m.CValue != "" {
			m.Change()
			menus = append(menus, m)
		}
	}
	return
}

func (d *Dao) Actives(c context.Context) (acs []*operate.Active, err error) {
	var rows *sql.Rows
	if rows, err = d.activeGet.Query(c); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		ac := &operate.Active{}
		if err = rows.Scan(&ac.ID, &ac.ParentID, &ac.Name, &ac.Background, &ac.Type, &ac.Content); err != nil {
			return
		}
		ac.Change()
		acs = append(acs, ac)
	}
	return
}
