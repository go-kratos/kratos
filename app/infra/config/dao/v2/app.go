package v2

import (
	"go-common/app/infra/config/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AppByTree get token by Name.
func (d *Dao) AppByTree(zone, env string, treeID int64) (app *model.App, err error) {
	app = &model.App{}
	row := d.DB.Select("id,token").Where("tree_id = ? AND env=? AND zone=?", treeID, env, zone).Model(&model.DBApp{}).Row()
	if err = row.Scan(&app.ID, &app.Token); err != nil {
		log.Error("AppByTree(%v) error(%v)", treeID, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

// AppsByNameEnv get token by Name.
func (d *Dao) AppsByNameEnv(name, env string) (apps []*model.DBApp, err error) {
	if err = d.DB.Where("env = ? and  name like ?", env, "%"+name).Find(&apps).Error; err != nil {
		log.Error("AppsByNameEnv(%v) error(%v)", name, env)
		return
	}
	if len(apps) == 0 {
		err = ecode.NothingFound
	}
	return
}

// AppGet ...
func (d *Dao) AppGet(zone, env, token string) (app *model.App, err error) {
	app = &model.App{}
	row := d.DB.Select("id,token,env,zone,tree_id").Where("token = ? AND env= ? AND zone= ?", token, env, zone).Model(&model.DBApp{}).Row()
	if err = row.Scan(&app.ID, &app.Token, &app.Env, &app.Zone, &app.TreeID); err != nil {
		log.Error("AppGet zone(%v) env(%v) token(%v) error(%v)", zone, env, token, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}
