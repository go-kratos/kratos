package wechat

import (
	"context"

	"go-common/app/tool/saga/model"
	"go-common/library/log"
)

// Changes changes structure
type Changes struct {
	Adds []*model.ContactInfo
	Upts []*model.ContactInfo
	Dels []*model.ContactInfo
}

// SyncContacts sync the contacts from wechat work
func (w *Wechat) SyncContacts(c context.Context) (err error) {
	if err = w.AnalysisContacts(c); err != nil {
		return
	}
	/*if err = w.UpdateVisible(c); err != nil {
		return
	}*/
	return
}

// AnalysisContacts analysis the contact difference and save them
func (w *Wechat) AnalysisContacts(c context.Context) (err error) {
	var (
		contactsInDB   []*model.ContactInfo
		wechatContacts []*model.ContactInfo
		changes        = &Changes{}
	)
	if contactsInDB, err = w.dao.ContactInfos(); err != nil {
		return
	}
	if wechatContacts, err = w.QueryWechatContacts(c); err != nil {
		return
	}

	if changes, err = w.diffChanges(wechatContacts, contactsInDB); err != nil {
		return
	}

	if err = w.saveChanges(changes); err != nil {
		return
	}
	return
}

// QueryWechatContacts query wechat contacts with access token
func (w *Wechat) QueryWechatContacts(c context.Context) (contacts []*model.ContactInfo, err error) {
	var (
		token string
	)
	if token, err = w.AccessToken(c, w.contact); err != nil {
		return
	}

	if contacts, err = w.dao.WechatContacts(c, token); err != nil {
		return
	}

	return
}

func (w *Wechat) saveChanges(changes *Changes) (err error) {
	var (
		contact *model.ContactInfo
	)

	log.Info("saveChanges add(%d), upt(%d), del(%d)", len(changes.Adds), len(changes.Upts), len(changes.Dels))
	for _, contact = range changes.Adds {
		if err = w.dao.CreateContact(contact); err != nil {
			return
		}
		log.Info("saveChanges add: %v", contact)
	}

	for _, contact = range changes.Upts {
		if err = w.dao.UptContact(contact); err != nil {
			return
		}
		log.Info("saveChanges upt: %v", contact)
	}

	for _, contact = range changes.Dels {
		if err = w.dao.DelContact(contact); err != nil {
			return
		}
		log.Info("saveChanges del: %v", contact)
	}

	return
}

func (w *Wechat) diffChanges(wechatContacts, contactsInDB []*model.ContactInfo) (changes *Changes, err error) {
	var (
		contact           *model.ContactInfo
		wechatContactsMap = make(map[string]*model.ContactInfo)
		contactsInDBMap   = make(map[string]*model.ContactInfo)
		wechatContactIDs  []string
		dbContactsIDs     []string
		userID            string
	)
	changes = new(Changes)
	for _, contact = range wechatContacts {
		wechatContactsMap[contact.UserID] = contact
		wechatContactIDs = append(wechatContactIDs, contact.UserID)
	}
	for _, contact = range contactsInDB {
		contactsInDBMap[contact.UserID] = contact
		dbContactsIDs = append(dbContactsIDs, contact.UserID)
	}

	// 分析变化
	for _, userID = range wechatContactIDs {
		contact = wechatContactsMap[userID]
		if w.inSlice(dbContactsIDs, userID) { // 企业微信联系人ID，在数据库中能找到
			if !contact.AlmostEqual(contactsInDBMap[userID]) { // 但是域不同
				contact.ID = contactsInDBMap[userID].ID
				changes.Upts = append(changes.Upts, contact)
			}
		} else {
			changes.Adds = append(changes.Adds, contact) // 这个联系人是新增的
		}
	}

	for _, userID = range dbContactsIDs {
		if !w.inSlice(wechatContactIDs, userID) {
			changes.Dels = append(changes.Dels, contactsInDBMap[userID])
		}
	}

	return
}

func (w *Wechat) inSlice(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// UpdateVisible update the visible property
func (w *Wechat) UpdateVisible(c context.Context) (err error) {
	var (
		user    *model.UserInfo
		users   []*model.UserInfo
		contact *model.ContactInfo
	)
	if users, err = w.querySagaVisible(c); err != nil {
		return
	}

	for _, user = range users {
		contact = &model.ContactInfo{UserID: user.UserID, VisibleSaga: true}
		if err = w.dao.UptContact(contact); err != nil {
			return
		}
	}
	return
}

func (w *Wechat) querySagaVisible(c context.Context) (users []*model.UserInfo, err error) {
	var (
		token string
	)
	if token, err = w.AccessToken(c, w.saga); err != nil {
		return
	}

	if users, err = w.dao.WechatSagaVisible(c, token, w.saga.AppID); err != nil {
		return
	}
	return
}
