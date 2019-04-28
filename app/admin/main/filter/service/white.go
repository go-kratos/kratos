package service

import (
	"context"

	"go-common/app/admin/main/filter/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
)

// AddAreaWhite .
func (s *Service) AddAreaWhite(c context.Context, content string, mode int8, areas []string, tps []int64, adid int64, name, comment string) (err error) {
	// contentID
	var (
		rule *model.WhiteInfo
		tx   *xsql.Tx
	)
	// check rule
	if err = s.checkArea(c, areas); err != nil {
		return
	}
	if err = s.checkReg(mode, content); err != nil {
		return
	}
	if err = s.checkBlackSample(mode, content); err != nil {
		return
	}
	if rule, err = s.dao.WhiteContent(c, content); err != nil {
		return
	}
	if rule != nil && rule.State == model.FilterStateNormal {
		return ecode.FilterDuplicateContent
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	var contentID int64
	if contentID, err = s.dao.UpsertWhiteContent(c, tx, content, comment, mode); err != nil {
		tx.Rollback()
		return
	}
	if rule != nil {
		contentID = rule.ID
	}
	if _, err = s.dao.DeleteAreaWhite(c, tx, contentID); err != nil {
		tx.Rollback()
		return
	}
	for _, area := range areas {
		for _, tp := range tps {
			if _, err = s.dao.UpsertAreaWhite(c, tx, area, tp, contentID); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	if _, err = s.dao.InsertWhiteLog(context.TODO(), tx, contentID, adid, name, comment, model.LogStateAdd); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	s.mission(func() {
		s.notifySearch(context.TODO(), areas)
	})
	return
}

// DeleteWhite .
func (s *Service) DeleteWhite(c context.Context, contentID, adid int64, name, reason string) (err error) {
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	if _, err = s.dao.DeleteWhiteContent(c, tx, contentID); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.DeleteAreaWhite(c, tx, contentID); err != nil {
		tx.Rollback()
		return
	}
	// log
	if _, err = s.dao.InsertWhiteLog(context.TODO(), tx, contentID, adid, name, reason, model.LogStateDel); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	s.mission(func() {
		s.notifySearch(context.TODO(), []string{"common"})
	})
	return
}

// SearchWhite .
func (s *Service) SearchWhite(c context.Context, content, area string, pn, ps int64) (rs []*model.WhiteInfo, total int64, err error) {
	if total, err = s.dao.CountWhiteContent(c, content, area); err != nil {
		return
	}
	rs, err = s.dao.SearchWhiteContent(c, content, area, (pn-1)*ps, ps)
	if len(rs) == 0 {
		rs = []*model.WhiteInfo{}
	}
	return
}

// WhiteInfo get white info (filter_white_content & filter_white_area) by content id
func (s *Service) WhiteInfo(c context.Context, id int64) (rule *model.WhiteInfo, err error) {
	return s.dao.WhiteInfo(c, id)
}

// EditWhite 白名单编辑不允许修改内容
func (s *Service) EditWhite(c context.Context, content string, mode int8, areas []string, tps []int64, adid int64, name, comment, reason string) (err error) {
	var (
		rule *model.WhiteInfo
		tx   *xsql.Tx
	)
	// check rule
	if err = s.checkArea(c, areas); err != nil {
		return
	}
	if err = s.checkReg(mode, content); err != nil {
		return
	}
	if err = s.checkBlackSample(mode, content); err != nil {
		return
	}
	if rule, err = s.dao.WhiteContent(c, content); err != nil {
		return
	}
	// 如果不存在该内容则报错
	if rule == nil || rule.State != model.FilterStateNormal {
		return ecode.RequestErr
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	if _, err = s.dao.UpdateWhiteContent(c, tx, mode, content, comment); err != nil {
		tx.Rollback()
		return
	}
	// 先删除旧有的
	if _, err = s.dao.DeleteAreaWhite(c, tx, rule.ID); err != nil {
		tx.Rollback()
		return
	}
	for _, area := range areas {
		for _, tp := range tps {
			if _, err = s.dao.UpsertAreaWhite(c, tx, area, tp, rule.ID); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	// log
	if _, err = s.dao.InsertWhiteLog(context.TODO(), tx, rule.ID, adid, name, reason, model.LogStateEdit); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	s.mission(func() {
		s.notifySearch(context.TODO(), areas)
	})
	return
}

// WhiteEditLog .
func (s *Service) WhiteEditLog(c context.Context, contentID int64) (ls []*model.Log, err error) {
	return s.dao.WhiteLog(c, contentID)
}
