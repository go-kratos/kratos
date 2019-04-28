package dao

import (
	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// FindUserByUserName find user by username.
func (d *Dao) FindUserByUserName(name string) (user *model.User, err error) {
	user = &model.User{}
	err = pkgerr.WithStack(d.db.Where("name = ?", name).First(user).Error)
	return
}

// FindUserByID find user by id.
func (d *Dao) FindUserByID(ID int64) (user *model.User, err error) {
	user = &model.User{}
	err = pkgerr.WithStack(d.db.Where("id = ?", ID).First(user).Error)
	return
}

// CreateUser create user.
func (d *Dao) CreateUser(user *model.User) (err error) {
	return pkgerr.WithStack(d.db.Create(user).Error)
}

// DelUser delete user.
func (d *Dao) DelUser(user *model.User) (err error) {
	return pkgerr.WithStack(d.db.Delete(user).Error)
}
