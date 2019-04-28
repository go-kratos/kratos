package service

import (
	"context"

	"go-common/app/interface/main/passport-login/model"
	"go-common/library/ecode"
)

var (
	tmpUserNotExist  = ecode.New(626)
	tmpDuplicateUser = ecode.New(652)
)

// CheckUser check username and pwd.
func (s *Service) CheckUser(c context.Context, username, pwd string) (decodeUser *model.DecodeUser, err error) {
	var (
		user *model.User
	)
	// tel
	if IsTel(username) {
		if user, err = s.telCheck(c, username); err != nil {
			return nil, err
		}
	} else if IsMail(username) { // mail
		if user, err = s.mailCheck(c, username); err != nil {
			return nil, err
		}
	} else { // userID
		if user, err = s.userIDCheck(c, username); err != nil {
			return nil, err
		}
	}
	decodeUser = user.Decode()
	// check use pwd .
	if err = s.checkUserPwd(pwd, decodeUser.Pwd, decodeUser.Salt); err != nil {
		decodeUser = nil
		return
	}
	decodeUser.Pwd = ""
	decodeUser.Salt = ""
	return
}

func (s *Service) userIDCheck(c context.Context, username string) (res *model.User, err error) {
	var (
		mid int64
	)
	if mid, err = s.dao.GetUserBaseByUserID(context.Background(), username); err != nil {
		return nil, err
	}
	if mid == 0 {
		return nil, tmpUserNotExist
	}
	if res, err = s.dao.GetUserByMid(context.Background(), mid); err != nil {
		return nil, err
	}
	if res == nil {
		return nil, tmpUserNotExist
	}
	return res, nil
}

func (s *Service) mailCheck(c context.Context, username string) (res *model.User, err error) {
	var (
		mid       int64
		midMail   int64
		midUserID int64
	)
	mail := s.doEncrypt(username)
	if midMail, err = s.dao.GetUserMailByMail(context.Background(), mail); err != nil {
		return nil, err
	}
	if midUserID, err = s.dao.GetUserBaseByUserID(context.Background(), username); err != nil {
		return nil, err
	}
	if midMail == 0 && midUserID == 0 {
		return nil, tmpUserNotExist
	}
	if midMail != 0 && midUserID != 0 {
		if midMail != midUserID {
			return nil, tmpDuplicateUser
		}
	}
	if midMail != 0 {
		mid = midMail
	}
	if midUserID != 0 {
		mid = midUserID
	}
	if res, err = s.dao.GetUserByMid(context.Background(), mid); err != nil {
		return
	}
	if res == nil {
		return nil, tmpUserNotExist
	}
	return res, nil
}

func (s *Service) telCheck(c context.Context, username string) (res *model.User, err error) {
	var (
		mid       int64
		midTel    int64
		midUserID int64
	)
	tel := s.doEncrypt(username)
	if midTel, err = s.dao.GetUserTelByTel(context.Background(), tel); err != nil {
		return nil, err
	}
	if midUserID, err = s.dao.GetUserBaseByUserID(context.Background(), username); err != nil {
		return nil, err
	}
	if midTel == 0 && midUserID == 0 {
		return nil, tmpUserNotExist
	}
	if midTel != 0 && midUserID != 0 {
		if midTel != midUserID {
			return nil, tmpDuplicateUser
		}
	}
	if midTel != 0 {
		mid = midTel
	}
	if midUserID != 0 {
		mid = midUserID
	}
	if res, err = s.dao.GetUserByMid(context.Background(), mid); err != nil {
		return nil, err
	}
	if res == nil {
		return nil, tmpUserNotExist
	}
	return res, nil
}
