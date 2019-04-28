package dao

import (
	"context"

	"go-common/app/job/main/workflow/model"
	"go-common/library/log"
)

// BusinessAttr .
func (d *Dao) BusinessAttr(c context.Context) (res []*model.BusinessAttr, err error) {
	if err = d.ReadORM.Table("workflow_business_attr").Select("id, bid, name, deal_type, expire_time, assign_type, assign_max, group_type").Find(&res).Error; err != nil {
		log.Error("d.BusinessAttr error(%v)", err)
	}
	return
}
