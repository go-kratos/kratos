package service

import (
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

// QueryUser query user info
func (s *Service) QueryUser(userName string) (user *model.User, err error) {
	return s.CreateUser(userName)
}

// CreateUser create user
func (s *Service) CreateUser(userName string) (userData *model.User, err error) {

	//此处因为业务因素，出现错误需要继续执行，不能retrun !!!
	if userData, err = s.dao.QueryUserByUserName(userName); err != nil {
		log.Error("s.dao.QueryUserByUserName err :(%v)", err)
	}
	if userData.ID == 0 {
		user := model.User{Name: userName, Email: userName + "@bilibili.com", Active: "1", Accept: -1}
		s.dao.AddUser(&user)
		userData, err = s.dao.QueryUserByUserName(userName)
	}
	return
}

// UpdateUser update user
func (s *Service) UpdateUser(user *model.User) error {
	return s.dao.UpdateUser(user)
}
