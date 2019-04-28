package service

import (
	"context"

	"go-common/app/service/main/member/model"
)

// NickUpdated get nickUpdated.
func (s *Service) NickUpdated(c context.Context, mid int64) (nickUpdated bool, err error) {
	if nickUpdated, err = s.mbDao.UserAttrDB(c, mid, model.NickUpdated); err != nil {
		return
	}
	return
}

// SetNickUpdated update isUpNickFree =1.
func (s *Service) SetNickUpdated(c context.Context, mid int64) (err error) {
	return s.mbDao.SetUserAttr(c, mid, model.NickUpdated)
}
