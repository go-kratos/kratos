package region

import (
	"context"

	"go-common/app/interface/main/app-intl/conf"
	"go-common/app/interface/main/app-intl/model/region"
	xsql "go-common/library/database/sql"
)

const (
	// region
	_secSQL = "SELECT rid,name,logo,rank,goto,param,plat,area,build,conditions FROM region WHERE state=1 AND reid!=0"
)

// Dao is.
type Dao struct {
	db  *xsql.DB
	get *xsql.Stmt
}

// New a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: xsql.NewMySQL(c.MySQL.Show),
	}
	// prepare
	d.get = d.db.Prepared(_secSQL)
	return
}

// Seconds get all second region.
func (d *Dao) Seconds(c context.Context) (rm map[int8]map[int16]*region.Region, err error) {
	rows, err := d.get.Query(c)
	if err != nil {
		return
	}
	defer rows.Close()
	rm = map[int8]map[int16]*region.Region{}
	for rows.Next() {
		a := &region.Region{}
		if err = rows.Scan(&a.Rid, &a.Name, &a.Logo, &a.Rank, &a.Goto, &a.Param, &a.Plat, &a.Area, &a.Build, &a.Condition); err != nil {
			return
		}
		if rs, ok := rm[a.Plat]; ok {
			rs[a.Rid] = a
		} else {
			rm[a.Plat] = map[int16]*region.Region{a.Rid: a}
		}
	}
	return
}

// Close close db resource.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
