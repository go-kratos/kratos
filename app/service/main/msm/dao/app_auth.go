package dao

import (
	"context"

	"go-common/app/service/main/msm/model"
	"go-common/library/log"
)

const (
	_allAppInfoSQL = "SELECT app_tree_id,app_id,`limit` FROM app"
	_allAppAuthSQL = "SELECT service_tree_id,app_tree_id,rpc_method,http_method,quota,mtime FROM app_auth"
)

// AllAppsInfo AllAppsInfo.
func (d *Dao) AllAppsInfo(c context.Context) (res map[int64]*model.AppInfo, err error) {
	rows, err := d.db.Query(c, _allAppInfoSQL)
	if err != nil {
		log.Error("d.apmDB.Query(app) error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.AppInfo)
	for rows.Next() {
		app := &model.AppInfo{}
		if err = rows.Scan(&app.AppTreeID, &app.AppID, &app.Limit); err != nil {
			log.Error("rows.Scan(app) error(%v)", err)
			return
		}
		res[app.AppTreeID] = app
	}
	return
}

// AllAppsAuth get all app auth info.
func (d *Dao) AllAppsAuth(c context.Context) (res map[int64]map[int64]*model.AppAuth, err error) {
	rows, err := d.db.Query(c, _allAppAuthSQL)
	if err != nil {
		log.Error("d.apmDB.Query(app_auth) error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]map[int64]*model.AppAuth)
	for rows.Next() {
		appAuth := &model.AppAuth{}
		if err = rows.Scan(&appAuth.ServiceTreeID, &appAuth.AppTreeID, &appAuth.RPCMethod, &appAuth.HTTPMethod, &appAuth.Quota, &appAuth.MTime); err != nil {
			log.Error("rows.Scan(appAuth) error(%v)", err)
			return
		}
		if _, b := res[appAuth.ServiceTreeID]; b {
			res[appAuth.ServiceTreeID][appAuth.AppTreeID] = appAuth
		} else {
			authMap := make(map[int64]*model.AppAuth)
			authMap[appAuth.AppTreeID] = appAuth
			res[appAuth.ServiceTreeID] = authMap
		}
	}
	return
}
