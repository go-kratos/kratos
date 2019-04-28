package service

import (
	"context"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddLabel  create new label
func (s *Service) AddLabel(label *model.Label) (err error) {
	if _, err = s.dao.QueryLabel(label); err == nil {
		err = ecode.MelloiLabelExistErr
		return
	}
	label.Active = 1
	return s.dao.AddLabel(label)
}

// QueryLabel query all labels
func (s *Service) QueryLabel(c context.Context) ([]*model.Label, error) {
	return s.dao.QueryLabels(c)
}

// DeleteLabel  delete label by id
func (s *Service) DeleteLabel(id int64) error {
	if id <= 0 {
		return ecode.RequestErr
	}
	return s.dao.DeleteLabel(id)
}

// AddLabelRelation create new label relation
func (s *Service) AddLabelRelation(lr *model.LabelRelation) (err error) {
	label := &model.Label{ID: lr.LabelID}
	if label, err = s.dao.QueryLabel(label); err != nil {
		return ecode.RequestErr
	}

	// 存在相同的记录
	if _, err = s.dao.QueryLabelExist(lr); err != nil {
		return ecode.MelloiLabelExistErr
	}

	// 每个脚本|任务，最多有2个label
	var lre []*model.LabelRelation
	relation := model.LabelRelation{Type: lr.Type, TargetID: lr.TargetID}
	if lre, err = s.dao.QueryLabelRelation(&relation); err != nil {
		return err
	}
	if len(lre) >= 2 {
		return ecode.MelloiLabelCountErr
	}

	lr.Description = label.Description
	lr.Color = label.Color
	lr.LabelName = label.Name
	lr.Active = 1
	if err = s.dao.AddLabelRelation(lr); err != nil {
		log.Error("s.dao.AddLabelRelation err :(%v)", err)
		return
	}
	return
}

// DeleteLabelRelation delete label relation by id
func (s *Service) DeleteLabelRelation(id int64) (err error) {
	if id <= 0 {
		return ecode.RequestErr
	}
	// 标签不存在
	var exist bool
	if exist, err = s.dao.CheckLabelRelationExist(id); err != nil {
		return
	}
	if !exist {
		return ecode.MelloiLabelRelationNotExist
	}

	return s.dao.DeleteLabelRelation(id)
}

// QueryLabelRelation query label relation by id , type,  targetid
func (s *Service) QueryLabelRelation(lre *model.LabelRelation) (lr []*model.LabelRelation, err error) {
	return s.dao.QueryLabelRelation(lre)
}

// QueryLabelRelationByIDs query label relation by ids
func (s *Service) QueryLabelRelationByIDs(ids []int64) (lr []*model.LabelRelation, err error) {
	return s.dao.QueryLabelRelationByIDs(ids)
}
