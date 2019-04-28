package gorm

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/aegis/model/business"
	"go-common/library/log"
)

// GetConfigs .
func (d *Dao) GetConfigs(c context.Context, bizid int64) (cfgs []*business.BizCFG, err error) {
	if err = d.orm.Table("business_config").Where("business_id=? AND state=0", bizid).Scan(&cfgs).Error; err != nil {
		log.Error("GetURL error(%v)", err)
	}
	return
}

// GetConfig .
func (d *Dao) GetConfig(c context.Context, bizid int64, tp int8) (config string, err error) {
	if err = d.orm.Table("business_config").Select("`config`").
		Where("business_id=? AND type=? AND state=0", bizid, tp).
		Order("mtime DESC").Limit(1).
		Row().Scan(&config); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("GetConfig error(%v)", err)
		return
	}
	return
}

// ActiveConfigs 所有任务配置
func (d *Dao) ActiveConfigs(c context.Context) (configs []*business.BizCFG, err error) {
	configs = []*business.BizCFG{}
	if err = d.orm.Where("state=0").Find(&configs).Error; err != nil {
		log.Error("ActiveConfigs find error(%v)", err)
	}
	return
}

// AddBizConfig 每个业务每种配置只有一条
func (d *Dao) AddBizConfig(c context.Context, cfg *business.BizCFG) (lastid int64, err error) {
	if err = d.orm.Table("business_config").Where("business_id=? AND `type`=?", cfg.BusinessID, cfg.TP).
		Assign(map[string]interface{}{
			"config": cfg.Config,
			"state":  cfg.State,
		}).FirstOrCreate(cfg).Error; err != nil {
		log.Error("AddBizConfig error(%v)", err)
	}

	lastid = cfg.ID
	return
}

// EditBizConfig .
func (d *Dao) EditBizConfig(c context.Context, cfg *business.BizCFG) (err error) {
	if err = d.orm.Table("business_config").Where("id=?", cfg.ID).
		Update(map[string]interface{}{
			"config": cfg.Config,
			"state":  cfg.State,
		}).Error; err != nil {
		log.Error("EditBizConfig error(%v)", err)
	}
	return
}
