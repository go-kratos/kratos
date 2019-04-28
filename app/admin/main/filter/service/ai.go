package service

import (
	"context"

	"go-common/app/admin/main/filter/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AiConfig get AI config.
func (s *Service) AiConfig(c context.Context) (res *model.AiConfig) {
	res = &model.AiConfig{}
	res.Threshold = s.conf.Ai.Threshold
	res.TrueScore = s.conf.Ai.TrueScore
	return
}

// AiWhite Ai white mids list.
func (s *Service) AiWhite(c context.Context, pn, ps int) (res []*model.AiWhite, total int64, err error) {
	if res, err = s.dao.AiWhite(c, pn, ps); err != nil {
		log.Error("AiWhite, s.dao.AiWhite(%d %d) error(%v)", pn, ps, err)
		return
	}
	if total, err = s.dao.AiWhiteCount(c); err != nil {
		log.Error("AiWhite, s.dao.AiWhiteCount() error(%v)", err)
	}
	return
}

// AiWhiteAdd AI white mid add.
func (s *Service) AiWhiteAdd(c context.Context, mid int64) (err error) {
	id, _ := s.dao.AiWhiteByMid(c, mid)
	if id > 0 {
		err = ecode.FilterInvalidAIWhiteUID
		return
	}
	if _, err = s.dao.InsertAiWhite(c, mid); err != nil {
		log.Error("AiWhiteAdd, s.dao.InsertAiWhite(%d) error(%v)", mid, err)
	}
	return
}

// AiWhiteEdit AI white mid edit.
func (s *Service) AiWhiteEdit(c context.Context, mid int64, state int8) (err error) {
	if _, err = s.dao.EditAiWhite(c, mid, state); err != nil {
		log.Error("AiWhiteEdit, s.dao.EditAiWhite(%d,%d) error(%v)", mid, state, err)
	}
	return
}

// AiScore get AI score.
func (s *Service) AiScore(c context.Context, content string) (res *model.AiScore, err error) {
	if res, err = s.dao.AiScore(c, content); err != nil {
		log.Error("AiScore, s.dao.AiScore(%s) error(%v)", content, err)
	}
	return
}

// AiCaseAdd AI case add.
func (s *Service) AiCaseAdd(c context.Context, aiCase *model.AiCase) (err error) {
	if _, err = s.dao.InsertAiCase(c, aiCase); err != nil {
		log.Error("AiCaseAdd, s.dao.AiCaseAdd(%+v) error(%v)", aiCase, err)
	}
	return
}
