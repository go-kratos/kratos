package service

import (
	"context"

	"go-common/app/admin/main/videoup-task/model"
)

func (s *Service) diffVideoOper(vp *model.VideoParam) (conts []string) {
	if vp.TagID > 0 {
		var operType int8
		if vp.Status >= model.VideoStatusOpen {
			operType = model.OperTypeOpenTag
		} else {
			operType = model.OperTypeRecicleTag
		}
		conts = append(conts, model.Operformat(operType, "tagid", vp.TagID, model.OperStyleTwo))
	}
	if vp.Reason != "" {
		conts = append(conts, model.Operformat(model.OperTypeAduitReason, "reason", vp.Reason, model.OperStyleTwo))
	}
	if vp.TaskID > 0 {
		conts = append(conts, model.Operformat(model.OperTypeTaskID, "task", vp.TaskID, model.OperStyleTwo))
	}
	return
}

func (s *Service) addVideoOper(c context.Context, oper *model.VideoOper) (err error) {
	/*
		if oldOper, _ := s.dao.VideoOper(c, oper.Vid); oldOper != nil && oldOper.LastID == 1 {
			oper.LastID = oldOper.ID
			s.dao.AddVideoOper(c, oper.Aid, oper.UID, oper.Vid, oper.Attribute, oper.Status, oper.LastID, oper.Content, oper.Remark)
			return
		}
	*/
	if lastID, _ := s.dao.AddVideoOper(c, oper.Aid, oper.UID, oper.Vid, oper.Attribute, oper.Status, oper.LastID, oper.Content, oper.Remark); lastID > 0 {
		s.dao.UpVideoOper(c, lastID, lastID)
		return
	}
	return
}
