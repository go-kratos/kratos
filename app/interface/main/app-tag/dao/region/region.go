package region

import (
	"context"

	"go-common/app/interface/main/app-tag/conf"
	"go-common/app/interface/main/app-tag/model/region"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// region
	_allSQL = "SELECT r.rid,r.reid,r.name,r.logo,r.rank,r.goto,r.param,r.plat,r.area,r.build,r.conditions,r.uri,r.is_logo,l.name FROM region AS r, language AS l WHERE r.state=1 AND l.id=r.lang_id ORDER BY r.rank DESC"
	//region android
	_regionPlatSQL = "SELECT r.rid,r.reid,r.name,r.logo,r.rank,r.goto,r.param,r.plat,r.area,r.build,r.conditions,l.name FROM region AS r, language AS l WHERE r.plat=0 AND r.state=1 AND l.id=r.lang_id ORDER BY r.rank DESC"
)

type Dao struct {
	db         *sql.DB
	get        *sql.Stmt
	regionPlat *sql.Stmt
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	// prepare
	d.get = d.db.Prepared(_allSQL)
	d.regionPlat = d.db.Prepared(_regionPlatSQL)
	return
}

// GetAll get all region.
func (d *Dao) All(ctx context.Context) ([]*region.Region, error) {
	rows, err := d.get.Query(ctx)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	apps := []*region.Region{}
	for rows.Next() {
		a := &region.Region{}
		if err = rows.Scan(&a.Rid, &a.Reid, &a.Name, &a.Logo, &a.Rank, &a.Goto, &a.Param, &a.Plat, &a.Area, &a.Build, &a.Condition, &a.URI, &a.Islogo, &a.Language); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, err
}

// RegionPlat get android
func (d *Dao) RegionPlat(ctx context.Context) ([]*region.Region, error) {
	rows, err := d.regionPlat.Query(ctx)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	apps := []*region.Region{}
	for rows.Next() {
		a := &region.Region{}
		if err = rows.Scan(&a.Rid, &a.Reid, &a.Name, &a.Logo, &a.Rank, &a.Goto, &a.Param, &a.Plat, &a.Area, &a.Build, &a.Condition, &a.Language); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, err
}

// resource DB ping
func (d *Dao) Ping(c context.Context) (err error) {
	return d.db.Ping(c)
}
