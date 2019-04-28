package service

import (
	"context"

	"go-common/app/admin/main/vip/model"
	"go-common/library/log"
)

// AllVersion all version.
func (s *Service) AllVersion(c context.Context) (res []*model.VipAppVersion, err error) {
	if res, err = s.dao.AllVersion(c); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

// UpdateVersion update version.
func (s *Service) UpdateVersion(c context.Context, v *model.VipAppVersion) (err error) {
	if _, err = s.dao.UpdateVersion(c, v); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}
