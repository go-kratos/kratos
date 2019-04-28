package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_vipConfDialog = "vip_conf_dialog"
)

// DialogAll .
func (d *Dao) DialogAll(c context.Context, appID, platform int64, status string) (res []*model.ConfDialog, err error) {
	db := d.vip.Table(_vipConfDialog)
	if appID != 0 {
		db = db.Where("app_id=?", appID)
	}
	if platform != 0 {
		db = db.Where("platform=?", platform)
	}
	if len(status) > 0 {
		curr := time.Now().Format("2006-01-02 15:04:05")
		//padding：待生效，active：已经生效，inactive：已经失效
		switch status {
		case "padding":
			db = db.Where("stage = true AND start_time>?", curr)
		case "active":
			db = db.Where("stage = true AND start_time<=? AND (end_time = '1970-01-01 08:00:00' OR end_time >?)", curr, curr)
		case "inactive":
			db = db.Where("stage = false OR (end_time > '1970-01-01 08:00:00' AND end_time < ?)", curr)
		default:
			log.Info("query all dialog.")
		}
	}
	if err := db.Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// DialogByID vip price config by id.
func (d *Dao) DialogByID(c context.Context, id int64) (dlg *model.ConfDialog, err error) {
	dlg = &model.ConfDialog{}
	if err := d.vip.Table(_vipConfDialog).Where("id=?", id).First(dlg).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		}
		return nil, err
	}
	return
}

// DialogBy vip price config by .
func (d *Dao) DialogBy(c context.Context, appID, platform int64, id int64) (res []*model.ConfDialog, err error) {
	if err := d.vip.Table(_vipConfDialog).Where("stage = true AND app_id=? AND platform=? AND id<>?", appID, platform, id).Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// DialogSave .
func (d *Dao) DialogSave(c context.Context, arg *model.ConfDialog) (eff int64, err error) {
	db := d.vip.Table(_vipConfDialog).Save(arg)
	if err = db.Error; err != nil {
		return
	}
	eff = db.RowsAffected
	return
}

// DialogEnable .
func (d *Dao) DialogEnable(c context.Context, arg *model.ConfDialog) (eff int64, err error) {
	stage := map[string]interface{}{
		"stage":    arg.Stage,
		"end_time": time.Now(),
		"operator": arg.Operator,
	}
	db := d.vip.Table(_vipConfDialog).Where("id=?", arg.ID).Updates(stage)
	if err = db.Error; err != nil {
		return
	}
	eff = db.RowsAffected
	return
}

// DialogDel delete vip price config by id.
func (d *Dao) DialogDel(c context.Context, id int64) (eff int64, err error) {
	db := d.vip.Table(_vipConfDialog).Where("id=?", id).Delete(&model.ConfDialog{})
	if err = db.Error; err != nil {
		return
	}
	eff = db.RowsAffected
	return
}

// CountDialogByPlatID count dialog by platform id .
func (d *Dao) CountDialogByPlatID(c context.Context, plat int64) (count int64, err error) {
	if err := d.vip.Table(_vipConfDialog).Where("platform=?", plat).Count(&count).Error; err != nil {
		return 0, err
	}
	return
}
