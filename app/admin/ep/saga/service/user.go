package service

import "go-common/app/admin/ep/saga/model"

// UserInfo get username and email.
func (s *Service) UserInfo(userName string) (userInfo *model.User) {

	userInfo = &model.User{
		Name:  userName,
		EMail: userName + "@bilibili.com",
	}
	return
}
