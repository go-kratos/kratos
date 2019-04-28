package dao

import (
	"context"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
)

const (
	_vipConfPlatform = "vip_platform_config"
)

// PlatformAll .
func (d *Dao) PlatformAll(c context.Context, order string) (res []*model.ConfPlatform, err error) {
	db := d.vip.Table(_vipConfPlatform)
	if err := db.Where("is_del=0").Order("id " + order).Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// PlatformByID vip platform config by id.
func (d *Dao) PlatformByID(c context.Context, id int64) (re *model.ConfPlatform, err error) {
	re = &model.ConfPlatform{}
	if err := d.vip.Table(_vipConfPlatform).Where("id=?", id).First(re).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		}
		return nil, err
	}
	return
}

// PlatformSave .
func (d *Dao) PlatformSave(c context.Context, arg *model.ConfPlatform) (eff int64, err error) {
	db := d.vip.Table(_vipConfPlatform).Omit("ctime").Save(arg)
	if err = db.Error; err != nil {
		return
	}
	eff = db.RowsAffected
	return
}

// PlatformEnable .
// func (d *Dao) PlatformEnable(c context.Context, arg *model.ConfPlatform) (eff int64, err error) {
// 	isDel := map[string]interface{}{
// 		"is_del":    arg.IsDel,
// 		"operator": arg.Operator,
// 	}
// 	db := d.vip.Table(_vipConfPlatform).Where("id=?", arg.ID).Updates(isDel)
// 	if err = db.Error; err != nil {
// 		return
// 	}
// 	eff = db.RowsAffected
// 	return
// }

// PlatformDel delete vip platform config by id.
func (d *Dao) PlatformDel(c context.Context, id int64, operator string) (eff int64, err error) {
	isDel := map[string]interface{}{
		"is_del":   1,
		"operator": operator,
	}
	db := d.vip.Table(_vipConfPlatform).Where("id=?", id).Updates(isDel)
	if err = db.Error; err != nil {
		return
	}
	eff = db.RowsAffected
	return
}

// PlatformTypes .
func (d *Dao) PlatformTypes(c context.Context) (res []*model.TypePlatform, err error) {
	db := d.vip.Table(_vipConfPlatform)
	if err := db.Select("id, platform_name").Where("is_del=0").Order("id").Find(&res).Error; err != nil {
		return nil, err
	}
	return
}
