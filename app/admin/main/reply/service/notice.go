package service

import (
	"context"
	"fmt"
	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ListNotice ListNotice
func (s *Service) ListNotice(c context.Context, page, pageSize int64) (nts []*model.Notice, total int64, err error) {
	offset := (page - 1) * pageSize
	total, err = s.dao.CountNotice(c)
	if err != nil {
		return
	}
	nts, err = s.dao.ListNotice(c, offset, pageSize)
	if err != nil {
		return
	}
	return nts, total, err
}

// GetNotice GetNotice
func (s *Service) GetNotice(c context.Context, id uint32) (nt *model.Notice, err error) {
	return s.dao.Notice(c, id)
}

// UpdateNotice UpdateNotice
func (s *Service) UpdateNotice(c context.Context, nt *model.Notice) (err error) {
	var ntGet *model.Notice
	ntGet, err = s.dao.Notice(c, nt.ID)
	if err != nil || ntGet == nil {
		err = ecode.NothingFound
		return
	}
	if ntGet.Status == model.StatusOnline {
		err = s.checkConflict(c, nt)
		if err != nil {
			return
		}
	}
	_, err = s.dao.UpdateNotice(c, nt)
	return
}

// CreateNotice CreateNotice
func (s *Service) CreateNotice(c context.Context, nt *model.Notice) (lastID int64, err error) {
	lastID, err = s.dao.CreateNotice(c, nt)
	if lastID <= 0 {
		log.Error("create notice failed!last_id not found")
		err = fmt.Errorf("create notice failed!last_id not found")
		return
	}
	return
}

// DeleteNotice DeleteNotice
func (s *Service) DeleteNotice(c context.Context, id uint32) (err error) {
	_, err = s.dao.DeleteNotice(c, id)
	return err
}

func (s *Service) checkConflict(c context.Context, nt *model.Notice) error {
	var nts []*model.Notice
	var err error
	nts, err = s.dao.RangeNotice(c, nt.Plat, nt.StartTime, nt.EndTime)
	if err != nil {
		return err
	}
	for _, data := range nts {
		//如果ID相同说明是自己
		if data.ID == nt.ID {
			continue
		}
		//如果为web平台则必然冲突
		if data.Plat == model.PlatWeb {
			return ecode.ReplyNoticeConflict
		}
		//如果客户端类型不同则跳过检查
		if data.ClientType != "" && nt.ClientType != "" && data.ClientType != nt.ClientType {
			continue
		}
		//为每一种版本情况检查
		if data.Condition == model.ConditionEQ {
			if nt.Condition == model.ConditionEQ && nt.Build == data.Build {
				return ecode.ReplyNoticeConflict
			}
			if nt.Condition == model.ConditionGT && nt.Build <= data.Build {
				return ecode.ReplyNoticeConflict
			}
			if nt.Condition == model.ConditionLT && nt.Build >= data.Build {
				return ecode.ReplyNoticeConflict
			}
		} else if data.Condition == model.ConditionGT {
			if nt.Condition == model.ConditionEQ {
				if nt.Build >= data.Build {
					return ecode.ReplyNoticeConflict
				}
			} else if nt.Condition == model.ConditionLT {
				if nt.Build >= data.Build {
					return ecode.ReplyNoticeConflict
				}
			} else {
				return ecode.ReplyNoticeConflict
			}
		} else if data.Condition == model.ConditionLT {
			if nt.Condition == model.ConditionEQ {
				if nt.Build <= data.Build {
					return ecode.ReplyNoticeConflict
				}
			} else if nt.Condition == model.ConditionGT {
				if nt.Build <= data.Build {
					return ecode.ReplyNoticeConflict
				}
			} else {
				return ecode.ReplyNoticeConflict
			}
		}
	}
	return nil
}

// UpdateNoticeStatus UpdateNoticeStatus
func (s *Service) UpdateNoticeStatus(c context.Context, status model.NoticeStatus, id uint32) (err error) {
	//检测客户端在同一时间内是否存在另外一条已发布的公告，如果存在则不允许发布
	if status == model.StatusOnline {
		var nt *model.Notice
		nt, err = s.dao.Notice(c, id)
		if err != nil || nt == nil {
			return
		}
		err = s.checkConflict(c, nt)
		if err != nil {
			return
		}
	}

	_, err = s.dao.UpdateNoticeStatus(c, status, id)
	return err
}
