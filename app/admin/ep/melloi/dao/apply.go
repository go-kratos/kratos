package dao

import (
	"go-common/app/admin/ep/melloi/model"
)

// QueryApply query apply list
func (d *Dao) QueryApply(apply *model.Apply, pn, ps int32) (qar *model.QueryApplyResponse, err error) {
	qar = &model.QueryApplyResponse{}
	err = d.DB.Table(model.Apply{}.TableName()).Where(model.Apply{
		ID: apply.ID, From: apply.From, To: apply.To, Status: apply.Status, Active: model.ApplyValid}).
		Count(&qar.TotalSize).Offset((pn - 1) * ps).Limit(ps).Order("id desc").Find(&qar.ApplyList).Error
	qar.PageSize = ps
	qar.PageNum = pn
	return
}

// QueryUserApplyList query user apply list
func (d *Dao) QueryUserApplyList(userName string) (applyList []*model.Apply, err error) {
	applyList = []*model.Apply{}
	err = d.DB.Table(model.Apply{}.TableName()).Where("`from`=?", userName).
		Where("`active`=?", 1).Where("`status`=?", model.ApplyValid).Find(&applyList).Error
	return
}

//QueryApplyByID query apply by id
func (d *Dao) QueryApplyByID(id int64) (apply *model.Apply, err error) {
	apply = &model.Apply{}
	err = d.DB.Table(model.Apply{}.TableName()).Where("id = ?", id).First(apply).Error
	return
}

//AddApply add apply
func (d *Dao) AddApply(apply *model.Apply) error {
	return d.DB.Model(&model.Apply{}).Create(apply).Error
}

// UpdateApply update apply info
func (d *Dao) UpdateApply(apply *model.Apply) error {
	return d.DB.Model(&model.Apply{}).Updates(apply).Where("ID=?", apply.ID).Error
}

// DeleteApply delete apply info
func (d *Dao) DeleteApply(id int64) error {
	return d.DB.Model(&model.Apply{}).Where("ID=?", id).Update("active", model.ApplyInvalid).Error
}
