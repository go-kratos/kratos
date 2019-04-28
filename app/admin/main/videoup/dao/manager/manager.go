package manager

import (
	"context"

	"go-common/app/admin/main/videoup/model/manager"
	"go-common/library/log"
)

const (
	_userInfoSQL = "SELECT id,username,nickname,state FROM user WHERE id = ?"
)

// User get manager user by id
func (d *Dao) User(c context.Context, id int64) (u *manager.User, err error) {
	var (
		row = d.managerDB.QueryRow(c, _userInfoSQL, id)
	)
	u = &manager.User{}
	if err = row.Scan(&u.ID, &u.UserName, &u.NickName, &u.State); err != nil {
		log.Error("row.Scan error(%v)", err)
		return
	}
	return
}
