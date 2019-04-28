package manager

import (
	"context"

	"go-common/app/job/main/videoup-report/model/manager"
	"go-common/library/log"
)

const (
	_userSQL = "SELECT u.id,u.username,d.name as department FROM user u left join user_department d on d.id=u.department_id where u.id=?"
)

// User get manager user profile
func (d *Dao) User(c context.Context, id int64) (user *manager.User, err error) {
	user = &manager.User{}
	if err = d.db.QueryRow(c, _userSQL, id).Scan(&user.ID, &user.Username, &user.Department); err != nil {
		log.Error("User db.row.Scan error(%v) uid(%d)", err, id)
	}
	return
}
