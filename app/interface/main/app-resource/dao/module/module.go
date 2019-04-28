package module

import (
	"context"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/module"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_moduleDevSQL = `SELECT p.id,d.name FROM resource_department AS d,resource_pool AS p WHERE d.id=p.department_id AND d.is_deleted=0 AND p.action=1 AND p.deleted=0`
	_moduleAllSQL = `SELECT rp.id,r.id,r.version,rp.name,rp.compresstype,rf.type,rf.url,rf.md5,rf.size,rf.file_type,rf.from_ver FROM resource_pool AS rp,resource AS r,resource_file AS rf 
	WHERE rp.id=r.pool_id AND r.id=rf.resource_id ORDER BY rp.id ASC,r.version DESC,rf.from_ver DESC`
	_resourceConfigSQL = `SELECT id,resource_id,stime,etime,valid,valid_test,default_package,is_wifi FROM resource_config WHERE is_deleted=0`
	_resourceLimitSQL  = `SELECT l.config_id,l.column,l.condition,l.value FROM resource_limit AS l WHERE l.is_deleted=0`
)

// Dao is a module dao.
type Dao struct {
	db *sql.DB
}

// New new a module dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	return
}

// ModuleDev get modelu dev
func (d *Dao) ModuleDev(ctx context.Context) (res map[string]*module.ResourcePool, err error) {
	rows, err := d.db.Query(ctx, _moduleDevSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[string]*module.ResourcePool)
	for rows.Next() {
		// NOTE: don't ask way the map key is "pink", stupid demand! fuck!!!
		// rp := &module.ResourcePool{ID: 1, Name: "pink"}
		r := &module.Resource{}
		if err = rows.Scan(&r.ID, &r.Name); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		var (
			resTmp *module.ResourcePool
			ok     bool
		)
		if resTmp, ok = res[r.Name]; !ok {
			resTmp = &module.ResourcePool{ID: r.ID, Name: r.Name}
			res[r.Name] = resTmp
		}
		resTmp.Resources = append(resTmp.Resources, r)
	}
	return
}

// ModuleAll get module
func (d *Dao) ModuleAll(ctx context.Context) (res map[int][]*module.Resource, err error) {
	rows, err := d.db.Query(ctx, _moduleAllSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int][]*module.Resource)
	for rows.Next() {
		r := &module.Resource{}
		if err = rows.Scan(&r.ID, &r.ResID, &r.Version, &r.Name, &r.Compresstype, &r.Type, &r.URL, &r.MD5, &r.Size, &r.Increment, &r.FromVer); err != nil {
			log.Error("d.db.rows.Scan error(%v)", err)
			return
		}
		r.URL = model.URLHTTPS(r.URL)
		res[r.ID] = append(res[r.ID], r)
	}
	return
}

// ResourceCondition
func (d *Dao) ResourceConfig(ctx context.Context) (conditions map[int]*module.Condition, err error) {
	rows, err := d.db.Query(ctx, _resourceConfigSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	conditions = make(map[int]*module.Condition)
	for rows.Next() {
		var (
			condition = &module.Condition{}
		)
		if err = rows.Scan(&condition.ID, &condition.ResID, &condition.STime, &condition.ETime, &condition.Valid, &condition.ValidTest, &condition.Default, &condition.IsWifi); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		conditions[condition.ResID] = condition
	}
	return
}

// ResourceLimit
func (d *Dao) ResourceLimit(ctx context.Context) (columns map[int]map[string][]*module.Column, err error) {
	rows, err := d.db.Query(ctx, _resourceLimitSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	columns = make(map[int]map[string][]*module.Column)
	for rows.Next() {
		var (
			cid         int
			column      = &module.Column{}
			columnValue string
			cm          map[string][]*module.Column
			ok          bool
		)
		if err = rows.Scan(&cid, &columnValue, &column.Condition, &column.Value); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if cm, ok = columns[cid]; !ok {
			cm = make(map[string][]*module.Column)
			columns[cid] = cm
		}
		cm[columnValue] = append(cm[columnValue], column)
	}
	return
}
