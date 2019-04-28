package egg

import (
	"context"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model/static"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_eggSQL = `SELECT e.id,ep.plat,ep.conditions,ep.build,ep.url,ep.md5,ep.size FROM egg AS e,egg_plat AS ep 
	WHERE e.id=ep.egg_id AND e.stime<? AND e.etime>? AND e.publish=1 AND e.delete=0 AND ep.deleted=0`
)

type Dao struct {
	db *xsql.DB
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: xsql.NewMySQL(c.MySQL.Show),
	}
	return
}

// Egg select all egg
func (d *Dao) Egg(ctx context.Context, now time.Time) (res map[int8][]*static.Static, err error) {
	res = map[int8][]*static.Static{}
	rows, err := d.db.Query(ctx, _eggSQL, now, now)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		s := &static.Static{}
		if err = rows.Scan(&s.Sid, &s.Plat, &s.Condition, &s.Build, &s.URL, &s.Hash, &s.Size); err != nil {
			log.Error("egg rows.Scan error(%v)", err)
			return
		}
		if s.URL == "" {
			continue
		}
		s.StaticChange()
		res[s.Plat] = append(res[s.Plat], s)
	}
	err = rows.Err()
	return
}
