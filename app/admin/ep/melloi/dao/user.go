package dao

import (
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

// AddUser add user by user object
func (d *Dao) AddUser(user *model.User) error {
	return d.DB.Create(user).Error
}

// QueryUserByUserName query user info by userName
func (d *Dao) QueryUserByUserName(userName string) (rows *model.User, err error) {
	rows = &model.User{}
	if err = d.DB.Where("name = ?", userName).First(rows).Error; err != nil {
		log.Error("d.db.Where error(%v)", err)
	}
	return
}

// UpdateUser update user info
func (d *Dao) UpdateUser(user *model.User) error {
	return d.DB.Model(&model.User{}).Update(user).Where("ID=?", user.ID).Error
}
