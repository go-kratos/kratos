package service

import (
	"context"
	"go-common/library/log"

	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"
)

// QueryUserInfo query user info.
func (s *Service) QueryUserInfo(c context.Context, username string) (userInfo *model.User, err error) {
	if userInfo, err = s.dao.FindUserByUserName(username); err != nil {
		return
	}

	if userInfo.ID == 0 {
		user := model.User{Name: username, EMail: username + "@bilibili.com"}
		s.dao.CreateUser(&user)
		userInfo, err = s.dao.FindUserByUserName(username)
	}

	return
}

// HttpSyncWechatContacts Http Sync Wechat Contacts.
func (s *Service) HttpSyncWechatContacts(ctx context.Context) (err error) {
	s.taskCache.Do(context.TODO(), func(ctx context.Context) {
		s.SyncWechatContacts()
	})
	return
}

// SyncWechatContacts Sync Wechat Contacts.
func (s *Service) SyncWechatContacts() (err error) {
	s.syncWechatContactsLock.Lock()
	defer s.syncWechatContactsLock.Unlock()
	log.Info("start")
	var (
		wechatContacts   []*model.WechatContact
		contactInfosInDB []*model.ContactInfo
	)

	if wechatContacts, err = s.dao.WechatContacts(context.Background()); err != nil {
		return
	}

	if contactInfosInDB, err = s.dao.QueryAllContactInfos(); err != nil {
		return
	}

	for _, wechatContact := range wechatContacts {
		var (
			tmpContactInfoInDB *model.ContactInfo
		)
		for _, contactInfoInDB := range contactInfosInDB {
			if contactInfoInDB.UserID == wechatContact.UserID {
				tmpContactInfoInDB = contactInfoInDB
				break
			}
		}
		if tmpContactInfoInDB != nil {
			//update
			tmpContactInfoInDB.UserName = wechatContact.EnglishName
			tmpContactInfoInDB.NickName = wechatContact.Name

			if err = s.dao.UpdateContactInfo(tmpContactInfoInDB); err != nil {
				continue
			}
		} else {
			//add
			tmpContactInfoInDB = &model.ContactInfo{
				UserName: wechatContact.EnglishName,
				NickName: wechatContact.Name,
				UserID:   wechatContact.UserID,
			}

			if err = s.dao.InsertContactInfo(tmpContactInfoInDB); err != nil {
				continue
			}
		}
	}
	log.Info("end")
	return
}

// AccessToBugly Access To Bugly.
func (s *Service) AccessToBugly(c context.Context, username string) (isAccess bool) {
	var (
		userInfo *model.User
		err      error
	)

	if userInfo, err = s.QueryUserInfo(c, username); err != nil {
		return
	}

	isAccess = userInfo.VisibleBugly

	return
}

// UpdateUserVisibleBugly Update User Visible Bugly.
func (s *Service) UpdateUserVisibleBugly(c context.Context, username, updateUsername string, visibleBugly bool) (err error) {
	var (
		hasRight bool
		userInfo *model.User
	)

	for _, super := range s.c.Bugly.SuperOwner {
		if username == super {
			hasRight = true
			break
		}
	}

	if !hasRight {
		err = ecode.AccessDenied
		return
	}

	if userInfo, err = s.dao.FindUserByUserName(updateUsername); err != nil {
		return
	}

	if userInfo.ID == 0 {
		err = ecode.NothingFound
		return
	}

	err = s.dao.UpdateUserVisibleBugly(userInfo.ID, visibleBugly)

	return
}
