package dao

import (
	"context"

	"go-common/app/admin/main/open/model"
)

// AddApp .
func (d *Dao) AddApp(c context.Context, g *model.App) error {
	return d.DB.Table("dm_apps").Create(g).Error
}

// DelApp .
func (d *Dao) DelApp(c context.Context, appid int64) error {
	return d.DB.Table("dm_apps").Where("appid = ?", appid).Update("enabled", 0).Error
}

// UpdateApp .
func (d *Dao) UpdateApp(c context.Context, arg *model.AppParams) error {
	return d.DB.Table("dm_apps").Where("appid = ?", arg.AppID).Update("app_name", arg.AppName).Error
}

// ListApp .
func (d *Dao) ListApp(c context.Context, t *model.AppListParams) (res []*model.App, err error) {
	db := d.DB.Table("dm_apps").Where("enabled = ?", 1)
	if t.AppKey != "" {
		db = db.Where("appkey = ?", t.AppKey)
	}
	if t.AppName != "" {
		db = db.Where("app_name = ?", t.AppName)
	}
	err = db.Find(&res).Error
	return
}
