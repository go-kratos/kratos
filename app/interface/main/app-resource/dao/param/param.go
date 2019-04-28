package param

import (
	"context"
	"fmt"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model/param"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// get param key and value
	_getAllSQL = "SELECT name,value,plat,build,conditions FROM param WHERE state=0 AND plat!=9"
)

// Dao is a param dao.
type Dao struct {
	db  *sql.DB
	get *sql.Stmt
}

// New new a param dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	// prepare
	d.get = d.db.Prepared(_getAllSQL)
	return
}

// All get all param
func (d *Dao) All(ctx context.Context) (m map[string][]*param.Param, err error) {
	rows, err := d.get.Query(ctx)
	if err != nil {
		log.Error("d.get error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	m = map[string][]*param.Param{}
	var _key = "param_%d"
	for rows.Next() {
		p := &param.Param{}
		if err = rows.Scan(&p.Name, &p.Value, &p.Plat, &p.Build, &p.Condition); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		p.Change()
		key := fmt.Sprintf(_key, p.Plat)
		m[key] = append(m[key], p)
	}
	return
}

// Close close memcache resource.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
