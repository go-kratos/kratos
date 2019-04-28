package plugin

import (
	"context"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model/plugin"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getSQL = "SELECT `name`,`package`,`policy`,`ver_code`,`ver_name`,`size`,`md5`,`url`,`enable`,`force`,`clear`,`min_build`,`max_build`,`base_code`,`base_name`,`desc`,`coverage` FROM plugin WHERE `enable`=1 AND `state`=0"
)

type Dao struct {
	db        *sql.DB
	pluginGet *sql.Stmt
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	// prepare
	d.pluginGet = d.db.Prepared(_getSQL)
	return
}

func (d *Dao) All(c context.Context) (psm map[string][]*plugin.Plugin, err error) {
	rows, err := d.pluginGet.Query(c)
	if err != nil {
		log.Error("query error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	psm = map[string][]*plugin.Plugin{}
	for rows.Next() {
		p := &plugin.Plugin{}
		if err = rows.Scan(&p.Name, &p.Package, &p.Policy, &p.VerCode, &p.VerName, &p.Size, &p.MD5, &p.URL, &p.Enable, &p.Force, &p.Clear, &p.MinBuild, &p.MaxBuild, &p.BaseCode, &p.BaseName, &p.Desc, &p.Coverage); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		if p.MaxBuild != 0 && p.MaxBuild < p.MinBuild {
			continue
		}
		psm[p.Name] = append(psm[p.Name], p)
	}
	return psm, err
}

// Close close memcache resource.
func (dao *Dao) Close() {
	if dao.db != nil {
		dao.db.Close()
	}
}

func (dao *Dao) PingDB(c context.Context) (err error) {
	return dao.db.Ping(c)
}
