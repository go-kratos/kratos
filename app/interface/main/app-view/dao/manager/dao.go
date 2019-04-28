package region

import (
	"context"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/manager"
	xsql "go-common/library/database/sql"
)

const (
	_relateSQL = "SELECT `id`,`param`,`goto`,`title`,`resource_ids`,`tag_ids`,`archive_ids`,`rec_reason`,`position`,`plat_ver`, `stime`,`etime` FROM app_rcmd_pos WHERE `state`=1"
)

type Dao struct {
	db  *xsql.DB
	get *xsql.Stmt
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: xsql.NewMySQL(c.MySQL.Show),
	}
	// prepare
	d.get = d.db.Prepared(_relateSQL)
	return
}

// Relate get all relate rec.
func (d *Dao) Relate(c context.Context) (rs []*manager.Relate, err error) {
	rows, err := d.get.Query(c)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &manager.Relate{}
		if err = rows.Scan(&r.ID, &r.Param, &r.Goto, &r.Title, &r.ResourceIDs, &r.TagIDs, &r.ArchiveIDs, &r.RecReason, &r.Position, &r.PlatVer, &r.STime, &r.ETime); err != nil {
			return
		}
		r.Change()
		rs = append(rs, r)
	}
	return
}

// Close close db resource.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
