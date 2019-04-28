package dao

import (
	"context"

	"go-common/app/job/main/workflow/model"
)

// ChallByIDs get chall list by ids.
func (d *Dao) ChallByIDs(c context.Context, cids []int64) (res map[int64]*model.Chall, err error) {
	if len(cids) <= 0 {
		return
	}
	res = make(map[int64]*model.Chall)
	cList := make([]*model.Chall, 0, len(cids))
	if err = d.ReadORM.Table("workflow_chall").Select("id, business, dispatch_state, dispatch_time").Where("id IN (?)", cids).Find(&cList).Error; err != nil {
		return
	}
	for _, c := range cList {
		res[c.ID] = c
	}
	return
}

// UpDispatchStateByIDs update by ids.
func (d *Dao) UpDispatchStateByIDs(c context.Context, cids []int64, dispatchState int64) (err error) {
	err = d.WriteORM.Table("workflow_chall").Where("id IN (?)", cids).Update("dispatch_state", dispatchState).Error
	return
}

// UpDispatchStateAdminIDByIds .
func (d *Dao) UpDispatchStateAdminIDByIds(c context.Context, cids []int64, dispatchState, assignAdminid int64) (err error) {
	err = d.WriteORM.Table("workflow_chall").Where("id IN (?)", cids).Update("dispatch_state", dispatchState).Update("assignee_adminid", assignAdminid).Error
	return
}
