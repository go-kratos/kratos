package service

import (
	"context"
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	notice "go-common/app/service/bbq/notice-service/api/v1"
)

// PushRegister .
func (s *Service) PushRegister(c context.Context, args *notice.UserPushDev) (response *v1.PushRegisterResponse, err error) {
	err = s.dao.PushLogin(c, args)
	response = &v1.PushRegisterResponse{}
	return
}

// PushLogout .
func (s *Service) PushLogout(c context.Context, args *notice.UserPushDev) (response *v1.PushRegisterResponse, err error) {
	err = s.dao.PushLogout(c, args)
	response = &v1.PushRegisterResponse{}
	return
}

// PushCallback .
func (s *Service) PushCallback(c context.Context, args *v1.PushCallbackRequest, mid int64, buvid string) (response *v1.PushCallbackResponse, err error) {
	s.dao.PushCallback(c, args.TID, args.NID, mid, buvid)
	return
}
