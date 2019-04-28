package splash

import (
	"context"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model/splash"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_actAllSQL = `SELECT i.splash_id,s.animate,s.duration,s.type,s.times,s.goto,s.param,s.skip,s.starttime,s.endtime,s.platform,i.url,i.hash,i.width,i.height,s.area,s.conditions,s.build,s.operate,s.no_preview FROM 
					splash AS s, splash_image AS i WHERE s.id=i.splash_id AND s.publish=1 AND i.type=3 AND s.platform!=0 AND s.type!=2 AND s.type!=4 AND s.state=0 ORDER BY s.starttime DESC`
	_actBirthSQL = `SELECT i.splash_id,s.animate,s.duration,s.type,s.times,s.goto,s.param,s.skip,s.starttime,s.endtime,s.platform,i.url,i.hash,i.width,i.height,s.area,s.conditions,s.build,s.operate FROM 
					splash AS s, splash_image AS i WHERE s.id=i.splash_id AND s.publish=1 AND i.type=3 AND s.platform!=0 AND s.type=2 AND s.state=0 ORDER BY s.starttime DESC`
	_actVipSQL = `SELECT i.splash_id,s.animate,s.duration,s.type,s.times,s.goto,s.param,s.skip,s.starttime,s.endtime,s.platform,i.url,i.hash,i.width,i.height,s.area,s.conditions,s.build,s.operate FROM 
					splash AS s, splash_image AS i WHERE s.id=i.splash_id AND s.publish=1 AND i.type=3 AND s.platform!=0 AND s.type=4 AND s.state=0 ORDER BY s.starttime DESC`
)

// Dao is splash dao.
type Dao struct {
	resdb *sql.DB
	// splash_active
	actAll   *sql.Stmt
	actBirth *sql.Stmt
	actVip   *sql.Stmt
}

// New new splash dao and return.
func New(c *conf.Config) *Dao {
	d := &Dao{
		resdb: sql.NewMySQL(c.MySQL.Resource),
	}
	// splah_active
	d.actAll = d.resdb.Prepared(_actAllSQL)
	d.actBirth = d.resdb.Prepared(_actBirthSQL)
	d.actVip = d.resdb.Prepared(_actVipSQL)
	return d
}

// GetActiveAll get all splash from table splash_active.
func (d *Dao) ActiveAll(ctx context.Context) (res []*splash.Splash, err error) {
	rows, err := d.actAll.Query(ctx)
	if err != nil {
		log.Error("dao.Exec(), err (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		sp := &splash.Splash{}
		if err = rows.Scan(&sp.ID, &sp.Animate, &sp.Duration, &sp.Type, &sp.Times, &sp.Goto, &sp.Param, &sp.Skip, &sp.Start,
			&sp.End, &sp.Plat, &sp.Image, &sp.Hash, &sp.Width, &sp.Height, &sp.Area, &sp.Condition, &sp.Build, &sp.Operate, &sp.NoPreview); err != nil {
			log.Error("rows.Scan err (%v)", err)
			res = nil
			return
		}
		sp.PlatChange()
		if sp.Operate == 1 {
			res = append([]*splash.Splash{sp}, res...)
		} else {
			res = append(res, sp)
		}
	}
	return
}

// ActiveBirth from table splash and splash_image.
func (d *Dao) ActiveBirth(ctx context.Context) (res []*splash.Splash, err error) {
	rows, err := d.actBirth.Query(ctx)
	if err != nil {
		log.Error("dao.Exec(), err (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		sp := &splash.Splash{}
		if err = rows.Scan(&sp.ID, &sp.Animate, &sp.Duration, &sp.Type, &sp.Times, &sp.Goto, &sp.Param, &sp.Skip, &sp.Start,
			&sp.End, &sp.Plat, &sp.Image, &sp.Hash, &sp.Width, &sp.Height, &sp.Area, &sp.Condition, &sp.Build, &sp.Operate); err != nil {
			log.Error("rows.Scan err (%v)", err)
			res = nil
			return
		}
		sp.PlatChange()
		sp.BirthDate()
		res = append(res, sp)
	}
	return
}

// ActiveVip form table vip splash
func (d *Dao) ActiveVip(ctx context.Context) (res []*splash.Splash, err error) {
	rows, err := d.actVip.Query(ctx)
	if err != nil {
		log.Error("dao.Exec(), err (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		sp := &splash.Splash{}
		if err = rows.Scan(&sp.ID, &sp.Animate, &sp.Duration, &sp.Type, &sp.Times, &sp.Goto, &sp.Param, &sp.Skip, &sp.Start,
			&sp.End, &sp.Plat, &sp.Image, &sp.Hash, &sp.Width, &sp.Height, &sp.Area, &sp.Condition, &sp.Build, &sp.Operate); err != nil {
			log.Error("rows.Scan err (%v)", err)
			res = nil
			return
		}
		sp.PlatChange()
		res = append(res, sp)
	}
	return
}

// Close close memcache resource.
func (dao *Dao) Close() {
	if dao.resdb != nil {
		dao.resdb.Close()
	}
}
