package service

import (
	"context"
	"fmt"

	"go-common/app/tool/saga/conf"
	"go-common/app/tool/saga/model"
	"go-common/app/tool/saga/service/mail"
	"go-common/app/tool/saga/service/wechat"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// CollectWachatUsers send required wechat visible users stored in memcache by email
func (s *Service) CollectWachatUsers(c context.Context) (err error) {
	var (
		contactInfo *model.ContactInfo
		userMap     = make(map[string]model.RequireVisibleUser)
		user        string
	)

	if err = s.d.RequireVisibleUsers(c, &userMap); err != nil {
		log.Error("get require visible user error(%v)", err)
		return
	}

	for k, v := range userMap {
		if contactInfo, err = s.d.QueryUserByID(k); err == nil {
			if contactInfo.VisibleSaga {
				continue
			}
		}
		user += v.UserName + " ,  " + v.NickName + "\n"
	}

	content := fmt.Sprintf("\n\n邮箱前缀     昵称\n\n%s", user)
	for _, addr := range conf.Conf.Property.ReportRequiredVisible.AlertAddrs {
		if err = mail.SendMail2(addr, "需添加的企业微信名单", content); err != nil {
			return
		}
	}
	if err = s.d.DeleteRequireVisibleUsers(c); err != nil {
		log.Error("Delete require visible user error(%v)", err)
		return
	}

	return
}

// SyncContacts sync the wechat contacts
func (s *Service) SyncContacts(c context.Context) (err error) {
	var (
		w = wechat.New(s.d)
	)

	if err = w.SyncContacts(c); err != nil {
		return
	}
	return
}

// synccontactsproc sync wechat contact procedure
func (s *Service) synccontactsproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("synccontactsproc panic(%v)", errors.WithStack(fmt.Errorf("%v", x)))
			go s.synccontactsproc()
			log.Info("synccontactsproc recover")
		}
	}()
	var err error
	if err = s.SyncContacts(context.TODO()); err != nil {
		log.Error("s.SyncContacts err (%+v)", err)
	}
}
