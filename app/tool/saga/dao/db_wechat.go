package dao

import (
	"regexp"
	"strings"

	"go-common/app/tool/saga/model"
	"go-common/library/log"

	"github.com/pkg/errors"
	pkgerr "github.com/pkg/errors"
)

var (
	regUserID = regexp.MustCompile(`^\d+$`)
)

// QueryUserByUserName query user by user name
func (d *Dao) QueryUserByUserName(userName string) (contactInfo *model.ContactInfo, err error) {
	contactInfo = &model.ContactInfo{}
	err = pkgerr.WithStack(d.mysql.Where(&model.ContactInfo{UserName: userName}).First(contactInfo).Error)
	return
}

// QueryUserByID query user by user ID
func (d *Dao) QueryUserByID(userID string) (contactInfo *model.ContactInfo, err error) {
	contactInfo = &model.ContactInfo{}
	err = pkgerr.WithStack(d.mysql.Where(&model.ContactInfo{UserID: userID}).First(contactInfo).Error)
	return
}

// UserIds query user ids for the user names
func (d *Dao) UserIds(userNames []string) (userIds string, err error) {
	var (
		userName    string
		ids         []string
		contactInfo *model.ContactInfo
	)

	if len(userNames) == 0 {
		err = errors.Errorf("UserIds: userNames is empty!")
		return
	}

	for _, userName = range userNames {
		if contactInfo, err = d.QueryUserByUserName(userName); err != nil {
			log.Error("UserIds: no such user (%s) in db, err (%s)", userName, err.Error())
		}

		log.Info("UserIds: username (%s), userid (%s)", userName, contactInfo.UserID)
		if contactInfo.UserID != "" && regUserID.MatchString(contactInfo.UserID) {
			ids = append(ids, contactInfo.UserID)
		}
	}

	if len(ids) > 0 {
		userIds = strings.Join(ids, "|")
		err = nil
	} else {
		err = errors.Wrapf(err, "UserIds: failed to find all the users in db, what a pity!")
	}

	return
}

// ContactInfos query all the records in contact_infos
func (d *Dao) ContactInfos() (contactInfos []*model.ContactInfo, err error) {
	err = pkgerr.WithStack(d.mysql.Find(&contactInfos).Error)
	return
}

// CreateContact create contact info record
func (d *Dao) CreateContact(contact *model.ContactInfo) (err error) {
	err = pkgerr.WithStack(d.mysql.Create(contact).Error)
	return
}

// DelContact delete the contact info with the specified UserID
func (d *Dao) DelContact(contact *model.ContactInfo) (err error) {
	err = pkgerr.WithStack(d.mysql.Delete(contact).Error)
	return
}

// UptContact update the contact information
func (d *Dao) UptContact(contact *model.ContactInfo) (err error) {
	err = pkgerr.WithStack(d.mysql.Model(&model.ContactInfo{}).Where(&model.ContactInfo{UserID: contact.UserID}).Updates(*contact).Error)
	return
}
