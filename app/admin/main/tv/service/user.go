package service

import (
	"context"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

//UserInfo select user info by mid
func (s *Service) UserInfo(c context.Context, mid int64) (userInfo *model.TvUserInfoResp, err error) {
	if userInfo, err = s.dao.GetByMId(c, mid); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error("UserInfo (%v) error(%v)", userInfo, err)
		return
	}

	return
}
