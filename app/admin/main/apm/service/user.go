package service

import (
	"context"

	"go-common/app/admin/main/apm/model/user"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// Permit check user permission.
func (s *Service) Permit(c context.Context, username string, rule string) (err error) {
	for _, man := range s.c.Superman {
		if username == man {
			return
		}
	}
	usr := &user.User{}
	if err = s.DB.Where("username=?", username).First(usr).Error; err != nil {
		log.Error("s.DB.User(%s) user error(%v)", username, err)
		err = ecode.AccessDenied
		return
	}
	if user.Rules[rule].Permit == user.PermitDefault {
		return
	}
	cnt := 0
	if err = s.DB.Model(&user.Rule{}).Where("user_id=? AND rule=?", usr.ID, rule).Count(&cnt).Error; err != nil {
		log.Error("s.DB.User(%s) count error(%v)", username, err)
		err = ecode.AccessDenied
		return
	}
	if cnt == 0 {
		log.Warn("s.DB.User(%s) count=0", username)
		err = ecode.AccessDenied
	}
	return
}

// GetDefaultPermission get the modules and rules which have default permission
func (s *Service) GetDefaultPermission(c context.Context) (modules []string, rules []string) {
	for m, mp := range user.Modules {
		if mp.Permit == user.PermitDefault {
			modules = append(modules, m)
		}
	}
	for r, rp := range user.Rules {
		if rp.Permit == user.PermitDefault {
			rules = append(rules, r)
		}
	}
	return
}

// GetUser get user info by username if it exists, otherwise create the user info
func (s *Service) GetUser(c context.Context, username string) (usr *user.User, err error) {
	usr = &user.User{}
	err = s.DB.Where("username = ?", username).First(usr).Error
	if err == gorm.ErrRecordNotFound {
		usr.UserName = username
		usr.NickName = username
		err = s.DB.Create(usr).Error
	}
	if err != nil {
		log.Error("apmSvc.GetUser error(%v)", err)
		return
	}
	s.ranksCache.Lock()
	if s.ranksCache.Map[username] != nil {
		usr.AvatarURL = s.ranksCache.Map[username].AvatarURL
	} else {
		usr.AvatarURL, _ = s.dao.GitLabFace(c, username)
	}
	s.ranksCache.Unlock()
	return
}
