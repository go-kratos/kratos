package dao

import (
	"context"
	"database/sql"

	"go-common/library/log"
)

const (
	_getUserRoleSQL = "SELECT  `role` FROM `auth_role` WHERE uid = ?"
)

// GetUserRole 用户角色
func (d *Dao) GetUserRole(c context.Context, uid int64) (role int8, err error) {
	err = d.mngDB.QueryRow(c, _getUserRoleSQL, uid).Scan(&role)
	if err != nil && err != sql.ErrNoRows {
		log.Error("d.managerDB.Query error(%v)", err)
		return
	}
	return role, nil
}
