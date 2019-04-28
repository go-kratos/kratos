package service

import (
	"context"

	"go-common/app/admin/ep/merlin/model"
)

// QueryUserInfo query user info.
func (s *Service) QueryUserInfo(c context.Context, username string) (userInfo *model.User, err error) {
	if userInfo, err = s.dao.FindUserByUserName(username); err != nil {
		user := model.User{Name: username, EMail: username + "@bilibili.com"}
		s.dao.CreateUser(&user)
		userInfo, err = s.dao.FindUserByUserName(username)
	}
	return
}
