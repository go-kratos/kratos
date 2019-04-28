package service

import (
	"context"
	"sort"

	"go-common/library/log"

	model "go-common/app/admin/main/macross/model/manager"
)

// User get user.
func (s *Service) User(c context.Context, system string) (res []*model.User) {
	users, ok := s.user[system]
	if ok {
		for _, user := range users {
			res = append(res, user)
		}
	}
	sort.Sort(model.Users(res))
	return
}

// SaveUser save user.
func (s *Service) SaveUser(c context.Context, roleID, userID int64, system, userName string) (err error) {
	var rows int64
	if userID == 0 {
		if rows, err = s.dao.AddUser(c, roleID, system, userName); err != nil {
			log.Error("s.dao.AddUser(%d, %s, %s) error(%v)", roleID, system, userName, err)
			return
		}
	} else {
		if rows, err = s.dao.UpUser(c, userID, roleID, userName); err != nil {
			log.Error("s.dao.UpUser(%d, %d, %s) error(%v)", userID, roleID, userName, err)
			return
		}
	}
	// update cache
	if rows != 0 {
		s.loadUserCache()
	}
	return
}

// DelUser del user.
func (s *Service) DelUser(c context.Context, userID int64) (err error) {
	var rows int64
	if rows, err = s.dao.DelUser(c, userID); err != nil {
		log.Error("s.dao.DelUser(%d) error(%v)", userID, err)
		return
	} else if rows != 0 {
		// update cache
		s.loadUserCache()
	}
	return
}
