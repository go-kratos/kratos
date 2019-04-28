package region

import (
	"context"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model/region"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// region
	_allSQL    = "SELECT r.rid,r.reid,r.name,r.logo,r.rank,r.goto,r.param,r.plat,r.area,r.build,r.conditions,r.uri,r.is_logo,r.type,r.is_rank,l.name FROM region AS r, language AS l WHERE r.state=1 AND l.id=r.lang_id ORDER BY r.rank DESC"
	_allSQL2   = "SELECT r.id,r.rid,r.reid,r.name,r.logo,r.rank,r.goto,r.param,r.plat,r.area,r.uri,r.is_logo,r.type,l.name FROM region_copy AS r, language AS l WHERE r.state=1 AND l.id=r.lang_id ORDER BY r.rank DESC"
	_limitSQL  = "SELECT l.id,l.rid,l.build,l.conditions FROM region_limit AS l,region_copy AS r WHERE l.rid=r.id"
	_configSQL = "SELECT c.id,c.rid,c.is_rank FROM region_rank_config AS c,region_copy AS r WHERE c.rid=r.id"
	//region android
	_regionPlatSQL = "SELECT r.rid,r.reid,r.name,r.logo,r.rank,r.goto,r.param,r.plat,r.area,l.name FROM region_copy AS r, language AS l WHERE r.plat=0 AND r.state=1 AND l.id=r.lang_id ORDER BY r.rank DESC"
)

type Dao struct {
	db         *sql.DB
	get        *sql.Stmt
	list       *sql.Stmt
	limit      *sql.Stmt
	config     *sql.Stmt
	regionPlat *sql.Stmt
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	// prepare
	d.get = d.db.Prepared(_allSQL)
	d.list = d.db.Prepared(_allSQL2)
	d.limit = d.db.Prepared(_limitSQL)
	d.config = d.db.Prepared(_configSQL)
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
		if err = rows.Scan(&a.Rid, &a.Reid, &a.Name, &a.Logo, &a.Rank, &a.Goto, &a.Param, &a.Plat, &a.Area, &a.Build, &a.Condition, &a.URI, &a.Islogo, &a.Rtype, &a.Entrance, &a.Language); err != nil {
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
		if err = rows.Scan(&a.Rid, &a.Reid, &a.Name, &a.Logo, &a.Rank, &a.Goto, &a.Param, &a.Plat, &a.Area, &a.Language); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, err
}

// AllList get all region.
func (d *Dao) AllList(ctx context.Context) ([]*region.Region, error) {
	rows, err := d.list.Query(ctx)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	apps := []*region.Region{}
	for rows.Next() {
		a := &region.Region{}
		if err = rows.Scan(&a.ID, &a.Rid, &a.Reid, &a.Name, &a.Logo, &a.Rank, &a.Goto, &a.Param, &a.Plat, &a.Area, &a.URI, &a.Islogo, &a.Rtype, &a.Language); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, err
}

// Limit region limits
func (d *Dao) Limit(ctx context.Context) (map[int64][]*region.Limit, error) {
	rows, err := d.limit.Query(ctx)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	limits := map[int64][]*region.Limit{}
	for rows.Next() {
		a := &region.Limit{}
		if err = rows.Scan(&a.ID, &a.Rid, &a.Build, &a.Condition); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return nil, err
		}
		limits[a.Rid] = append(limits[a.Rid], a)
	}
	return limits, err
}

// Config region configs
func (d *Dao) Config(ctx context.Context) (map[int64][]*region.Config, error) {
	rows, err := d.config.Query(ctx)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	configs := map[int64][]*region.Config{}
	for rows.Next() {
		a := &region.Config{}
		if err = rows.Scan(&a.ID, &a.Rid, &a.ScenesID); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return nil, err
		}
		a.ConfigChange()
		configs[a.Rid] = append(configs[a.Rid], a)
	}
	return configs, err
}

// Close close memcache resource.
func (dao *Dao) Close() {
	if dao.db != nil {
		dao.db.Close()
	}
}
