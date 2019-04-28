package dao

import (
	"context"

	"go-common/app/admin/ep/melloi/model"
)

// AddLabel  add new label
func (d *Dao) AddLabel(label *model.Label) error {
	return d.DB.Table(model.Label{}.LabelName()).Create(label).Error
}

// QueryLabels  query all labels
func (d *Dao) QueryLabels(c context.Context) (labels []*model.Label, err error) {
	err = d.DB.Table(model.Label{}.LabelName()).Where("active = ?", 1).Find(&labels).Error
	return
}

// QueryLabel query label by label name and id
func (d *Dao) QueryLabel(lb *model.Label) (label *model.Label, err error) {
	label = &model.Label{}
	err = d.DB.Table(model.Label{}.LabelName()).Where("active = ?", 1).
		Where(model.Label{Name: lb.Name, ID: lb.ID}).First(label).Error
	return
}

// DeleteLabel delete label
func (d *Dao) DeleteLabel(id int64) error {
	return d.DB.Table(model.Label{}.LabelName()).Where("id = ?", id).Update("active", 0).Error
}

// AddLabelRelation add label relation of target
func (d *Dao) AddLabelRelation(relation *model.LabelRelation) error {
	return d.DB.Table(model.LabelRelation{}.LabelRelationName()).Create(relation).Error
}

// QueryLabelRelation query label relation
func (d *Dao) QueryLabelRelation(lre *model.LabelRelation) (lr []*model.LabelRelation, err error) {
	err = d.DB.Table(model.LabelRelation{}.LabelRelationName()).
		Where("active = ?", 1).
		Where(model.LabelRelation{Type: lre.Type, LabelID: lre.LabelID, TargetID: lre.TargetID}).Find(&lr).Error
	return
}

// QueryLabelRelationByIDs  Query label relation by ids
func (d *Dao) QueryLabelRelationByIDs(ids []int64) (lr []*model.LabelRelation, err error) {
	err = d.DB.Table(model.LabelRelation{}.LabelRelationName()).Where(" active = ? ", 1).
		Where("id in (?)", ids).Find(lr).Error
	return
}

// CheckLabelRelationExist check label relation exist
func (d *Dao) CheckLabelRelationExist(id int64) (result bool, err error) {
	result = false
	lr := &model.LabelRelation{}
	err = d.DB.Table(model.LabelRelation{}.LabelRelationName()).Where("active = ?", 1).Where(" id = ?", id).First(lr).Error
	if lr.ID > 0 {
		result = true
	}
	return
}

// QueryLabelExist check label exist
func (d *Dao) QueryLabelExist(lre *model.LabelRelation) (lr *model.LabelRelation, err error) {
	lr = &model.LabelRelation{}
	err = d.DB.Table(model.LabelRelation{}.LabelRelationName()).
		Where("active = ?", 1).
		Where(model.LabelRelation{Type: lre.Type, LabelID: lre.LabelID, TargetID: lre.TargetID}).First(lr).Error
	return
}

// DeleteLabelRelation delete relation of label
func (d *Dao) DeleteLabelRelation(id int64) (err error) {
	return d.DB.Table(model.LabelRelation{}.LabelRelationName()).Where(" id = ?", id).Update("active", 0).Error
}
