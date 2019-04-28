package version

import (
	"context"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/version"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// version
	_allSQL = "SELECT id,plat,description,version,build,ptime FROM version WHERE state=0 ORDER BY ptime DESC"
	// version_update AND version
	_allUpsSQL = `SELECT up.id,up.channel,up.url,up.size,up.md5,v.build,v.version,v.description,up.sdkint_list,up.model,up.policy,up.is_force,
				up.is_push,up.policy_name,v.state,up.mtime,v.plat,up.policy_url,up.buvid_start,up.buvid_end
					FROM version v RIGHT JOIN version_update up ON v.id=up.vid AND v.plat in(0,7,8,9,11,12) WHERE up.state=1 AND v.state!=1 AND up.url!='' ORDER BY v.build DESC`
	_allUpsLimitSQL = `SELECT up_id,value,condi FROM version_update_limit ORDER BY value ASC`
	_allSosSQL      = `SELECT id,package,name,description,clear,ver_code,ver_name,url,size,enable_state,force_state,md5,min_build,coverage,sdkint,model 
					FROM version_so WHERE enable_state!=2 ORDER BY id DESC`
	_allIncreSQL = `SELECT i.id,i.target_version,i.target_build,i.target_id,i.source_version,i.source_build,i.source_id,i.task_id,i.file_path,i.url,i.md5,i.size,v.build,v.plat
						FROM version v RIGHT JOIN version_incremental i ON v.id=i.vid AND v.plat in(0,7,8,9,11,12) WHERE i.state=1 AND v.state!=1 AND i.url!='' ORDER BY v.build DESC`
	_allRnSQL = `SELECT id,deployment_key,bundle_id,url,md5,size,version FROM version_rn WHERE state=1 ORDER BY ctime DESC`
)

// Dao is a version dao
type Dao struct {
	db     *sql.DB
	get    *sql.Stmt
	update *sql.Stmt
	so     *sql.Stmt
	incre  *sql.Stmt
	rn     *sql.Stmt
	limit  *sql.Stmt
}

// New new a version dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	// prepare
	d.get = d.db.Prepared(_allSQL)
	d.update = d.db.Prepared(_allUpsSQL)
	d.so = d.db.Prepared(_allSosSQL)
	d.incre = d.db.Prepared(_allIncreSQL)
	d.rn = d.db.Prepared(_allRnSQL)
	d.limit = d.db.Prepared(_allUpsLimitSQL)
	return
}

// GetVersion get version.
func (d *Dao) All(ctx context.Context) (ver map[int8][]*version.Version, err error) {
	rows, err := d.get.Query(ctx)
	if err != nil {
		log.Error("d.get error(%v)", err)
		return
	}
	defer rows.Close()
	ver = map[int8][]*version.Version{}
	for rows.Next() {
		a := &version.Version{}
		if err = rows.Scan(&a.Id, &a.Plat, &a.Desc, &a.Version, &a.Build, &a.PTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		ver[a.Plat] = append(ver[a.Plat], a)
	}
	return
}

// GetUpdate get active version AND version_update
func (d *Dao) Updates(ctx context.Context) (ver map[int8]map[string][]*version.VersionUpdate, err error) {
	var (
		masterChannel  = "master"
		defaultChannel = "bili"
	)
	rows, err := d.update.Query(ctx)
	if err != nil {
		log.Error("d.update error(%v)", err)
		return
	}
	defer rows.Close()
	ver = map[int8]map[string][]*version.VersionUpdate{}
	for rows.Next() {
		a := &version.VersionUpdate{}
		if err = rows.Scan(&a.Id, &a.Channel, &a.Url, &a.Size, &a.MD5, &a.Build, &a.Version, &a.Desc, &a.SdkInts, &a.Model, &a.Policy, &a.IsForce,
			&a.IsPush, &a.PolicyName, &a.State, &a.Mtime, &a.Plat, &a.PolicyURL, &a.BuvidStart, &a.BuvidEnd); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		a.VersionUpdateChange()
		if a.State == 2 {
			a.IsGray = 1
		} else {
			a.IsGray = 0
		}
		if a.Channel == masterChannel && a.Plat != model.PlatAndroidTVYST {
			a.Channel = defaultChannel
		}
		if v, ok := ver[a.Plat]; !ok {
			ver[a.Plat] = map[string][]*version.VersionUpdate{
				a.Channel: []*version.VersionUpdate{a},
			}
		} else {
			v[a.Channel] = append(v[a.Channel], a)
		}
	}
	return
}

// Limits version update limits
func (d *Dao) Limits(ctx context.Context) (res map[int][]*version.UpdateLimit, err error) {
	res = map[int][]*version.UpdateLimit{}
	rows, err := d.limit.Query(ctx)
	if err != nil {
		log.Error("d.update error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		l := &version.UpdateLimit{}
		if err = rows.Scan(&l.ID, &l.BuildLimit, &l.Conditions); err != nil {
			log.Error("version update limit sql error(%v)", err)
			return
		}
		res[l.ID] = append(res[l.ID], l)
	}
	return
}

// GetSos get version_so
func (d *Dao) Sos(ctx context.Context) (ver map[string][]*version.VersionSo, err error) {
	rows, err := d.so.Query(ctx)
	if err != nil {
		log.Error("d.so error(%v)", err)
		return
	}
	defer rows.Close()
	ver = map[string][]*version.VersionSo{}
	for rows.Next() {
		s := &version.VersionSo{}
		if err = rows.Scan(&s.Id, &s.Package, &s.Name, &s.Description, &s.Clear, &s.Ver_code, &s.Ver_name, &s.Url, &s.Size, &s.Enable_state,
			&s.Force_state, &s.Md5, &s.Min_build, &s.Coverage, &s.Sdkint, &s.Model); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		ver[s.Name] = append(ver[s.Name], s)
	}
	return
}

// Incrementals version incrementals
func (d *Dao) Incrementals(ctx context.Context) (ver map[int8]map[string][]*version.Incremental, err error) {
	rows, err := d.incre.Query(ctx)
	if err != nil {
		log.Error("d.incre error(%v)", err)
		return
	}
	defer rows.Close()
	ver = map[int8]map[string][]*version.Incremental{}
	for rows.Next() {
		v := &version.Incremental{}
		if err = rows.Scan(&v.ID, &v.TargetVersion, &v.TargetBuild, &v.TargetID, &v.SourceVersion, &v.SourceBuild, &v.SourceID, &v.TaskID, &v.FilePath, &v.URL, &v.Md5,
			&v.Size, &v.Build, &v.Plat); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		if vi, ok := ver[v.Plat]; !ok {
			ver[v.Plat] = map[string][]*version.Incremental{
				v.SourceID: []*version.Incremental{v},
			}
		} else {
			vi[v.SourceID] = append(vi[v.SourceID], v)
		}
	}
	return
}

// Rn version rn
func (d *Dao) Rn(ctx context.Context) (ver map[string]map[string]*version.Rn, err error) {
	rows, err := d.rn.Query(ctx)
	if err != nil {
		log.Error("d.rn error(%v)", err)
		return
	}
	defer rows.Close()
	ver = map[string]map[string]*version.Rn{}
	for rows.Next() {
		v := &version.Rn{}
		if err = rows.Scan(&v.ID, &v.DeploymentKey, &v.BundleID, &v.URL, &v.Md5, &v.Size, &v.Version); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		if vi, ok := ver[v.DeploymentKey]; !ok {
			ver[v.DeploymentKey] = map[string]*version.Rn{
				v.Version: v,
			}
		} else if _, ok := vi[v.Version]; !ok {
			vi[v.Version] = v
		}
	}
	return
}

// Close close db resource.
func (dao *Dao) Close() {
	if dao.db != nil {
		dao.db.Close()
	}
}
