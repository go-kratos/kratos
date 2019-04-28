package dao

import (
	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// InsertContactInfo Insert Contact Info.
func (d *Dao) InsertContactInfo(contactInfo *model.ContactInfo) error {
	return pkgerr.WithStack(d.db.Create(contactInfo).Error)
}

// UpdateContactInfo Update Contact Info.
func (d *Dao) UpdateContactInfo(contactInfo *model.ContactInfo) error {
	return pkgerr.WithStack(d.db.Save(&contactInfo).Error)
}

// QueryContactInfoByUserID Query Contact Info By User ID
func (d *Dao) QueryContactInfoByUserID(userID string) (contactInfo *model.ContactInfo, err error) {
	contactInfo = &model.ContactInfo{}
	if err = d.db.Where("user_id = ?", userID).First(contactInfo).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// QueryContactInfoByUsername Query Contact Info By Username
func (d *Dao) QueryContactInfoByUsername(username string) (contactInfo *model.ContactInfo, err error) {
	contactInfo = &model.ContactInfo{}
	if err = d.db.Where("username = ?", username).First(contactInfo).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// QueryAllContactInfos Query All Contact Infos
func (d *Dao) QueryAllContactInfos() (contactInfos []*model.ContactInfo, err error) {
	err = pkgerr.WithStack(d.db.Find(&contactInfos).Error)
	return
}
