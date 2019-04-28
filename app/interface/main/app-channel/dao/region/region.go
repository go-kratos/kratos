package region

import (
	"context"

	"go-common/app/interface/main/app-channel/conf"
	"go-common/app/interface/main/app-channel/model/channel"
	"go-common/library/database/sql"
)

const (
	_allSQL2   = "SELECT r.id,r.rid,r.reid,r.name,r.logo,r.param,r.plat,r.area,r.uri,r.type,l.name FROM region_copy AS r, language AS l WHERE r.state=1 AND l.id=r.lang_id ORDER BY r.rank DESC"
	_limitSQL  = "SELECT l.id,l.rid,l.build,l.conditions FROM region_limit AS l,region_copy AS r WHERE l.rid=r.id"
	_configSQL = "SELECT c.id,c.rid,c.is_rank FROM region_rank_config AS c,region_copy AS r WHERE c.rid=r.id"
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

// AllList get all region.
func (d *Dao) AllList(ctx context.Context) (apps []*channel.Region, err error) {
	rows, err := d.db.Query(ctx, _allSQL2)
	if err != nil {
		return
	}
	defer rows.Close()
	apps = []*channel.Region{}
	for rows.Next() {
		a := &channel.Region{}
		if err = rows.Scan(&a.ID, &a.RID, &a.ReID, &a.Name, &a.Logo, &a.Param, &a.Plat, &a.Area, &a.URI, &a.Type, &a.Language); err != nil {
			return
		}
		apps = append(apps, a)
	}
	return
}

// Limit region limits
func (d *Dao) Limit(ctx context.Context) (limits map[int64][]*channel.RegionLimit, err error) {
	rows, err := d.db.Query(ctx, _limitSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	limits = map[int64][]*channel.RegionLimit{}
	for rows.Next() {
		a := &channel.RegionLimit{}
		if err = rows.Scan(&a.ID, &a.Rid, &a.Build, &a.Condition); err != nil {
			return
		}
		limits[a.Rid] = append(limits[a.Rid], a)
	}
	return
}

// Config region configs
func (d *Dao) Config(ctx context.Context) (configs map[int64][]*channel.RegionConfig, err error) {
	rows, err := d.db.Query(ctx, _configSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	configs = map[int64][]*channel.RegionConfig{}
	for rows.Next() {
		a := &channel.RegionConfig{}
		if err = rows.Scan(&a.ID, &a.Rid, &a.ScenesID); err != nil {
			return
		}
		configs[a.Rid] = append(configs[a.Rid], a)
	}
	return
}

// Close close memcache resource.
func (dao *Dao) Close() {
	if dao.db != nil {
		dao.db.Close()
	}
}
