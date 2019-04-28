package gorm

import (
	"context"

	taskmod "go-common/app/admin/main/aegis/model/task"
	"go-common/library/log"
)

// AddConfig add config
func (d *Dao) AddConfig(c context.Context, config *taskmod.Config, confJSON interface{}) (err error) {
	if config.ConfType == taskmod.TaskConfigRangeWeight { //粉丝数，等待时长，分组权重配置去重
		name := "%" + confJSON.(*taskmod.RangeWeightConfig).Name + "%"
		err = d.orm.Table("task_config").Where("business_id=? AND flow_id=? AND conf_type=? AND conf_json LIKE ?", config.BusinessID, config.FlowID, taskmod.TaskConfigRangeWeight, name).
			Assign(map[string]interface{}{
				"conf_json":   config.ConfJSON,
				"btime":       config.Btime,
				"etime":       config.Etime,
				"uid":         config.UID,
				"uname":       config.Uname,
				"description": config.Description,
			}).FirstOrCreate(config).Error
		return err
	}
	return d.orm.Create(config).Error
}

// UpdateConfig update config
func (d *Dao) UpdateConfig(c context.Context, config *taskmod.Config) (err error) {
	return d.orm.Model(&taskmod.Config{}).Where("id=?", config.ID).Update(config).Error
}

// SetStateConfig update config
func (d *Dao) SetStateConfig(c context.Context, id int64, state int8) (err error) {
	return d.orm.Model(&taskmod.Config{}).Where("id=?", id).Update("state", state).Error
}

// QueryConfigs list config
func (d *Dao) QueryConfigs(c context.Context, queryParams *taskmod.QueryParams) (configs []*taskmod.Config, count int64, err error) {
	db := d.orm.Model(&taskmod.Config{}).Where("conf_type=?", queryParams.ConfType).Where("state=?", queryParams.State)
	if queryParams.BusinessID > 0 {
		db = db.Where("business_id=?", queryParams.BusinessID)
	}
	if queryParams.FlowID > 0 {
		db = db.Where("flow_id=?", queryParams.FlowID)
	}
	if len(queryParams.Btime) > 0 && len(queryParams.Etime) > 0 {
		db = db.Where("mtime>=? AND mtime<=?", queryParams.Btime, queryParams.Etime)
	}
	if len(queryParams.ConfName) > 0 {
		db = db.Where("conf_json LIKE '%" + queryParams.ConfName + "%'")
	}

	if err = db.Count(&count).Offset((queryParams.Pn - 1) * queryParams.Ps).Order("mtime DESC").Limit(queryParams.Ps).Find(&configs).Error; err != nil {
		log.Error("query error(%v)", err)
		return
	}
	return
}

// DeleteConfig delete config
func (d *Dao) DeleteConfig(c context.Context, id int64) (err error) {
	return d.orm.Where("id=?", id).Delete(&taskmod.Config{}).Error
}
