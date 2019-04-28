package service

import (
	"context"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// PlatformAll .
func (s *Service) PlatformAll(c context.Context, order string) (res []*model.ConfPlatform, err error) {
	if res, err = s.dao.PlatformAll(c, order); err != nil {
		return
	}
	return
}

// PlatformByID .
func (s *Service) PlatformByID(c context.Context, arg *model.ArgID) (dlg *model.ConfPlatform, err error) {
	return s.dao.PlatformByID(c, arg.ID)
}

// PlatformSave .
func (s *Service) PlatformSave(c context.Context, arg *model.ConfPlatform) (eff int64, err error) {
	return s.dao.PlatformSave(c, arg)
}

// PlatformDel .
func (s *Service) PlatformDel(c context.Context, arg *model.ArgID, operator string) (eff int64, err error) {
	pcount, err := s.dao.CountVipPriceConfigByPlat(c, arg.ID)
	if err != nil {
		return
	}
	dcount, err := s.dao.CountDialogByPlatID(c, arg.ID)
	if err != nil {
		return
	}
	if pcount > 0 || dcount > 0 {
		err = ecode.VipPlatformConfDelErr
		return
	}
	eff, err = s.dao.PlatformDel(c, arg.ID, operator)
	log.Warn("user(%s) delete dialog(%d) effect row(%d)", operator, arg.ID, eff)
	return
}
