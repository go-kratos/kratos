package dao

import (
	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// FindUserByUserName find user by username.
func (d *Dao) FindUserByUserName(name string) (user *model.User, err error) {
	user = &model.User{}
	if err = d.db.Where("name = ?", name).First(user).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// FindUserByID find user by id.
func (d *Dao) FindUserByID(ID int64) (user *model.User, err error) {
	user = &model.User{}
	err = pkgerr.WithStack(d.db.Where("id = ?", ID).First(user).Error)
	return
}

// CreateUser create user.
func (d *Dao) CreateUser(user *model.User) error {
	return pkgerr.WithStack(d.db.Create(user).Error)
}

// DelUser delete user.
func (d *Dao) DelUser(user *model.User) error {
	return pkgerr.WithStack(d.db.Delete(user).Error)
}

// UpdateUserVisibleBugly Update User Visible Bugly.
func (d *Dao) UpdateUserVisibleBugly(ID int64, visibleBugly bool) error {
	return pkgerr.WithStack(d.db.Model(&model.User{}).Where("id = ?", ID).Update("visible_bugly", visibleBugly).Error)
}
